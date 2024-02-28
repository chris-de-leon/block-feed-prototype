import { MySql2Database } from "drizzle-orm/mysql2"
import * as testcontainers from "testcontainers"
import { drizzle } from "drizzle-orm/mysql2"
import { db } from "@api/shared/database"
import { utils } from "@api/shared/utils"
import { VerboseOptions } from "./utils"
import * as mysql from "mysql2/promise"
import * as path from "node:path"

export const DB_CONSTANTS = {
  MYSQL: {
    VERSION: "8.3.0",
    PORT: 3306,
    DB: "test",
  },
  ROLES: {
    ROOT: {
      UNAME: "root",
      PWORD: "password",
    },
    API: {
      UNAME: "api_user",
      PWORD: "password",
    },
  },
}

export const spawn = async (verbose: VerboseOptions) => {
  const rootDir = utils.getRootDir()

  const container = await testcontainers.GenericContainer.fromDockerfile(
    path.join(rootDir, "vendor", "mysql"),
    "Dockerfile",
  )
    .withBuildArgs({ MYSQL_VERSION: DB_CONSTANTS.MYSQL.VERSION })
    .build(`mysql-dev:${DB_CONSTANTS.MYSQL.PORT}`)

  return await container
    .withExposedPorts(DB_CONSTANTS.MYSQL.PORT)
    .withEnvironment({
      MYSQL_ROOT_PASSWORD: DB_CONSTANTS.ROLES.ROOT.PWORD,
      MYSQL_DATABASE: DB_CONSTANTS.MYSQL.DB,
    })
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

export const withDatabaseConn = async <T>(
  url: string,
  cb: (
    ctx: Readonly<{
      conn: MySql2Database<typeof db.schema>
    }>,
  ) => Promise<T> | T,
) => {
  const conn = await mysql.createConnection(url)
  try {
    return await cb({
      conn: drizzle(conn, {
        schema: db.schema,
        mode: "default",
      }),
    })
  } finally {
    await conn.end()
  }
}

export const getMySqlUrl = (
  container: testcontainers.StartedTestContainer,
  uname: string,
  pword: string,
) => {
  return `mysql://${uname}:${pword}@host.docker.internal:${container.getMappedPort(
    DB_CONSTANTS.MYSQL.PORT,
  )}/${DB_CONSTANTS.MYSQL.DB}`
}

export const getRootRoleUrl = (
  container: testcontainers.StartedTestContainer,
) =>
  getMySqlUrl(
    container,
    DB_CONSTANTS.ROLES.ROOT.UNAME,
    DB_CONSTANTS.ROLES.ROOT.PWORD,
  )

export const getApiRoleUrl = (container: testcontainers.StartedTestContainer) =>
  getMySqlUrl(
    container,
    DB_CONSTANTS.ROLES.API.UNAME,
    DB_CONSTANTS.ROLES.API.PWORD,
  )
