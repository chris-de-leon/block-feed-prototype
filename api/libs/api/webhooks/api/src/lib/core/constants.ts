import { database } from "@api/shared/database"
import { auth0 } from "@api/shared/auth0"

export const ENV_FILES = ["auth0.env", "node.env", "db.env"]

export const NAMESPACE = "webhooks"

export type Context = Readonly<{
  database: ReturnType<typeof database.core.createClient>
  auth0: ReturnType<typeof auth0.createClient>
}>

export const CONSTANTS = {
  MAX_BLOCKS: {
    MIN: 1,
    MAX: 10,
  },
  MAX_RETRIES: {
    MIN: 0,
    MAX: 10,
  },
  TIMEOUT_MS: {
    MIN: 0,
    MAX: 10000,
  },
  BLOCKCHAIN_ID: {
    MIN: 0,
    MAX: 1024,
  },
  LIMIT: {
    MIN: 0,
    MAX: 25,
  },
  OFFSET: {
    MIN: 0,
  },
}

export const OPERATIONS = {
  CREATE: (() => {
    const name = "Create" as const
    return {
      ID: `${NAMESPACE}${name}`,
      METHOD: "POST",
      NAME: name,
      PATH: `/${NAMESPACE}.${name}`,
    } as const
  })(),
  FIND_MANY: (() => {
    const name = "FindMany" as const
    return {
      ID: `${NAMESPACE}${name}`,
      METHOD: "GET",
      NAME: name,
      PATH: `/${NAMESPACE}.${name}`,
    } as const
  })(),
  FIND_ONE: (() => {
    const name = "FindOne" as const
    return {
      ID: `${NAMESPACE}${name}`,
      METHOD: "GET",
      NAME: name,
      PATH: `/${NAMESPACE}.${name}`,
    } as const
  })(),
  UPDATE: (() => {
    const name = "Update" as const
    return {
      ID: `${NAMESPACE}${name}`,
      METHOD: "POST",
      NAME: name,
      PATH: `/${NAMESPACE}.${name}`,
    } as const
  })(),
  REMOVE: (() => {
    const name = "Remove" as const
    return {
      ID: `${NAMESPACE}${name}`,
      METHOD: "POST",
      NAME: name,
      PATH: `/${NAMESPACE}.${name}`,
    } as const
  })(),
  DEPLOY: (() => {
    const name = "Deploy" as const
    return {
      ID: `${NAMESPACE}${name}`,
      METHOD: "POST",
      NAME: name,
      PATH: `/${NAMESPACE}.${name}`,
    } as const
  })(),
} as const