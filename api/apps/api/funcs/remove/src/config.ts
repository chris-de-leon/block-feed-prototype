import { funcsAPI } from "@api/api/funcs/api"
import { AWS } from "@serverless/typescript"
import { utils } from "@api/shared/utils"
import * as path from "path"

export const config: AWS["functions"] = {
  [funcsAPI.OPERATIONS.REMOVE.ID]: {
    environment: utils.resolveEnvVars(funcsAPI.ENV_FILES),
    logRetentionInDays: 1,
    handler: path.join(
      path.dirname(path.relative(process.cwd(), __filename)),
      "main.handler"
    ),
    events: [
      {
        http: {
          operationId: funcsAPI.OPERATIONS.REMOVE.NAME,
          path: funcsAPI.OPERATIONS.REMOVE.PATH,
          method: funcsAPI.OPERATIONS.REMOVE.METHOD,
        },
      },
    ],
  },
}