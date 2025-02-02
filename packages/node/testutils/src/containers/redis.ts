import { DOCKER_HOST, VerboseOptions } from "./utils"
import * as testcontainers from "testcontainers"

export const REDIS_CONSTANTS = {
  VERSION: "7.2.1-alpine3.18",
  PORT: 6379,
}

export const spawn = async (verbose: VerboseOptions) => {
  return await new testcontainers.GenericContainer(
    `redis:${REDIS_CONSTANTS.VERSION}`,
  )
    .withExposedPorts(REDIS_CONSTANTS.PORT)
    .withWaitStrategy(testcontainers.Wait.forListeningPorts())
    .withLogConsumer((stream) => {
      if (verbose.data) {
        stream.on("data", console.log)
      }
      if (verbose.errors) {
        stream.on("err", console.error)
      }
    })
    .start()
}

export const getRedisUrl = (container: testcontainers.StartedTestContainer) => {
  return `redis://${DOCKER_HOST}:${container.getMappedPort(
    REDIS_CONSTANTS.PORT,
  )}`
}
