import { RelayersAPI } from "@api/api/relayers/api"
import { AWS } from "@serverless/typescript"
import { utils } from "@api/shared/utils"
import * as path from "path"

export const config: AWS["functions"] = {
  [RelayersAPI.OPERATIONS.FIND_ONE.ID]: {
    environment: utils.resolveEnvVars(RelayersAPI.ENV_FILES),
    logRetentionInDays: 1,
    handler: path.join(
      path.dirname(path.relative(process.cwd(), __filename)),
      "main.handler",
    ),
    events: [
      {
        http: {
          operationId: RelayersAPI.OPERATIONS.FIND_ONE.NAME,
          path: RelayersAPI.OPERATIONS.FIND_ONE.PATH,
          method: RelayersAPI.OPERATIONS.FIND_ONE.METHOD,
        },
      },
    ],
  },
}