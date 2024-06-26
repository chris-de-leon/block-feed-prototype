import { gqlBadRequestError } from "../../../graphql/errors"
import { GraphQLAuthContext } from "../../../graphql/types"
import * as schema from "@block-feed/node-db"
import { and, eq } from "drizzle-orm"
import { z } from "zod"

export const zInput = z.object({
  id: z.string().uuid(),
})

export const handler = async (
  args: z.infer<typeof zInput>,
  ctx: GraphQLAuthContext,
) => {
  return await ctx.providers.mysql.drizzle.query.webhook
    .findFirst({
      where: and(
        eq(schema.webhook.customerId, ctx.clerk.user.id),
        eq(schema.webhook.id, args.id),
      ),
    })
    .then((result) => {
      if (result == null) {
        throw gqlBadRequestError(`record with id "${args.id}" does not exist`)
      }
      return result
    })
}
