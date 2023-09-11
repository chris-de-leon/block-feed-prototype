export const NAMESPACE = "funcs"

export const ENV_FILES = ["auth0.env", "db.env"]

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
} as const