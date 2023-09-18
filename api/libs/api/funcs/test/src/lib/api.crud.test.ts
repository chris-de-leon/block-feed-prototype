import { after, before, describe, it } from "node:test"
import { blockgateway } from "@api/block-gateway"
import { database } from "@api/shared/database"
import { testutils } from "@api/shared/testing"
import { auth0 } from "@api/shared/auth0"
import { randomUUID } from "node:crypto"
import { flow } from "@api/shared/flow"
import assert from "node:assert"

describe("Funcs CRUD Test", () => {
  const blockchain = new blockgateway.blockchains.FlowBlockchain(
    flow.createClient()
  )
  const api = testutils.getApi()
  const db = database.core.createClient()
  const a0 = auth0.createClient()
  const chainInfo = blockchain.getInfo()

  let headers = {}
  let auth0User: Awaited<ReturnType<typeof testutils.createAuth0User>> | null =
    null

  before(async () => {
    await testutils.wipeDB(db, database.schema.blockFeed.schemaName)

    const user = await testutils.createAuth0User(a0)
    const grnt = await user.getGrant()

    const height = await blockchain.getLatestBlockHeight()
    await database.queries.blockCursor.create(db, {
      blockchain: chainInfo.name,
      id: chainInfo.id,
      networkURL: chainInfo.networkURL,
      height: height,
    })

    auth0User = user
    headers = {
      Authorization: `bearer ${grnt.access_token}`,
    }
  })

  after(async () => {
    await auth0User?.cleanUp()
  })

  it("creates a function", async () => {
    await api
      .funcsCreate(
        {
          name: randomUUID(),
          cursorId: chainInfo.id,
        },
        { headers }
      )
      .then((result) => {
        assert.equal(result.data.count, 1)
      })

    await it("finds many functions", async () => {
      const firstId = await api
        .funcsFindMany(undefined, undefined, { headers })
        .then((result) => {
          assert.equal(result.data.length > 0, true)
          return result.data[0].id
        })

      await it("updates a function", async () => {
        const newName = "new-name"
        await api
          .funcsUpdate(
            {
              id: firstId,
              name: newName,
            },
            {
              headers,
            }
          )
          .then((result) => {
            assert.equal(result.data.count, 1)
          })

        await it("finds a function by ID", async () => {
          await api.funcsFindOne(firstId, { headers }).then((result) => {
            assert.equal(result.data.name, newName)
          })
        })

        await it("removes a function", async () => {
          await api.funcsRemove({ id: firstId }, { headers }).then((result) => {
            assert.equal(result.data.count, 1)
          })
        })
      })
    })
  })
})
