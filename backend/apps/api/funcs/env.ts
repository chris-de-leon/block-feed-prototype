import { utils } from "../../../libs/shared/utils"
import * as path from "path"

export const resolveEnvVars = () => {
  const appEnv = utils.getAppEnv()

  // NOTE: more env files can be defined below
  const envfiles: string[] = []
  envfiles.push("auth0.env")
  envfiles.push("db.env")

  const envRoot = path.join(utils.getRootDir(), "env", appEnv)
  const envPths = envfiles.map((filename) => path.join(envRoot, filename))
  return utils.getEnvVars(envPths)
}
