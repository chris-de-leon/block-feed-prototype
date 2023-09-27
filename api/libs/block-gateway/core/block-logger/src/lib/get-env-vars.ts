import { utils } from "@api/shared/utils"

export const getEnvVars = () => {
  const ENV_KEYS = {
    BLOCK_LOGGER_REDIS_URL: "BLOCK_LOGGER_REDIS_URL",
  } as const

  return {
    [ENV_KEYS.BLOCK_LOGGER_REDIS_URL]: new URL(
      utils.getRequiredEnvVar(ENV_KEYS.BLOCK_LOGGER_REDIS_URL)
    ),
  }
}
