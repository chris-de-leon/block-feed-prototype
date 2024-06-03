import { StripeEventFixture } from "./fixtures/stripe.event"
import { after, before, describe, it } from "node:test"
import * as testutils from "@block-feed/test-utils"
import * as assert from "node:assert"
import Stripe from "stripe"
import {
  DatabaseVendor,
  StripeVendor,
  RedisVendor,
  stripe,
  redis,
  db,
} from "@block-feed/vendors"
import {
  StripeWebhookEventConsumer,
  StripeWebhookEventProducer,
  RedisCacheFactory,
} from "../../src"

describe("Stripe Webhook Event Tests", () => {
  const testCleaner = new testutils.TestCleaner()
  const consumerGroupName = "stripe-webhook-event-consumers"
  const fakeStripeApiKey = "fake api key"
  const consumerName = "replica-1"
  const maxBlockMs = 100
  const fakeStripeEvents = [
    // NOTE: none of these events will trigger calls to the Stripe API
    StripeEventFixture("checkout.session.completed"),
    StripeEventFixture("customer.subscription.paused"),
    StripeEventFixture("customer.subscription.created"),
    StripeEventFixture("customer.subscription.deleted"),
    StripeEventFixture("customer.subscription.updated"),
    StripeEventFixture("customer.subscription.resumed"),
    { type: "this is an event with an unhandled/invalid type" },
    {
      // valid event type, bad event payload
      type: "customer.subscription.deleted",
      data: { object: { metadata: {} } },
    },
  ]
  const verbose = {
    database: true,
    container: {
      errors: true,
      data: false,
    },
  }

  let redisStreamC: testutils.containers.StartedTestContainer
  let redisCacheC: testutils.containers.StartedTestContainer
  let databaseC: testutils.containers.StartedTestContainer

  let redisStreamVendor: RedisVendor
  let redisCacheVendor: RedisVendor
  let stripeVendor: StripeVendor
  let dbVendor: DatabaseVendor

  before(async () => {
    try {
      // Creates containers
      process.stdout.write("Starting containers... ")
      ;[redisCacheC, redisStreamC, databaseC] = await Promise.all([
        testutils.containers.redis.spawn(verbose.container),
        testutils.containers.redis.spawn(verbose.container),
        testutils.containers.db.spawn(verbose.container),
      ])
      console.log("done!")

      // Schedule the containers for cleanup
      testCleaner.add(
        () => redisStreamC.stop(),
        () => redisCacheC.stop(),
        () => databaseC.stop(),
      )

      // Assigns database urls
      const redisStreamUrl =
        testutils.containers.redis.getRedisUrl(redisStreamC)
      const redisCacheUrl = testutils.containers.redis.getRedisUrl(redisCacheC)
      const apiDbUrl = testutils.containers.db.getApiUserUrl(databaseC)

      // Creates a stripe API client
      stripeVendor = stripe.client.create({
        STRIPE_API_KEY: fakeStripeApiKey,
      })

      // Creates a redis stream client
      redisStreamVendor = redis.client.create({
        REDIS_URL: redisStreamUrl,
      })

      // Creates a redis cache client
      redisCacheVendor = redis.client.create({
        REDIS_URL: redisCacheUrl,
      })

      // Creates a database vendor
      dbVendor = db.client.create({
        DB_LOGGING: verbose.database,
        DB_URL: apiDbUrl,
      })

      // Schedule the context for cleanup
      testCleaner.add(
        () => redisStreamVendor.client.quit(),
        () => redisCacheVendor.client.quit(),
        () =>
          new Promise((res, rej) => {
            dbVendor.pool.end((err) => {
              if (err != null) {
                rej(err)
              }
              res(null)
            })
          }),
      )
    } catch (err) {
      // The node test runner won't log any errors that occur in
      // the before hook, so we need to log them manually
      console.error(err)
      throw err
    }
  })

  after(async () => {
    await testCleaner.cleanUp(console.error)
  })

  it("Integration Test", { timeout: 5000 }, async () => {
    // Define the producer and consumer
    const producer = new StripeWebhookEventProducer(redisStreamVendor)
    const consumer = await StripeWebhookEventConsumer.build(
      RedisCacheFactory.createCheckoutSessionCache(
        stripeVendor,
        redisCacheVendor,
      ),
      stripeVendor,
      redisStreamVendor,
      dbVendor,
      consumerGroupName,
      consumerName,
    )

    // Produce some fake messages
    await Promise.allSettled(
      fakeStripeEvents.map((ev) => producer.produce(ev as any)),
    ).then((results) => {
      const errs = new Array<Error>()
      results.forEach((result) => {
        if (result.status === "rejected") {
          errs.push(new Error(String(result.reason)))
        } else {
          console.log(`Successfully added message with ID "${result.value}"`)
        }
      })
      if (errs.length > 0) {
        throw new Error(JSON.stringify(errs, null, 2))
      }
    })

    // Consume the messages
    const controller = new AbortController()
    const eventLog = new Map<
      string,
      { id: string; data: Stripe.Event; err: Error | undefined }
    >()
    for await (const events of consumer.consume(controller, maxBlockMs)) {
      if (events.length === 0) {
        controller.abort() // Exit if a block timeout occurs
      } else {
        events.forEach((ev) => eventLog.set(ev.id, ev))
      }
    }

    // Log the messages
    console.log(
      "Event log:\n",
      JSON.stringify(
        Object.fromEntries(eventLog.entries()),
        (_, v) => (v instanceof Error ? v.message : v),
        2,
      ),
    )

    // Verify that all the produced messages have been received
    assert.equal(eventLog.size, fakeStripeEvents.length)

    // Filter out the messages that have associated errors
    const errLogs = Array.from(eventLog.values())
      .flat()
      .filter((ev) => ev.err != null)

    // Only one message should have an error
    assert.equal(errLogs.length, 1)

    // Define expected error messages
    const expectedErrorStrings = {
      1: '{\n    "code": "invalid_type",\n    "expected": "string",\n    "received": "undefined",\n    "path": [\n      "userId"\n    ],\n    "message": "Required"\n  }',
      2: 'an error occurred while processing event "customer.subscription.deleted"',
    }

    // Validate the error message
    const errMsg = errLogs.at(0)?.err?.message
    assert.ok(
      errMsg?.includes(expectedErrorStrings[1]) &&
        errMsg?.includes(expectedErrorStrings[2]),
    )
  })
})
