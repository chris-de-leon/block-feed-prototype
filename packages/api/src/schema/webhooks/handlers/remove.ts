import { GraphQLAuthContext } from "../../../graphql/types"
import { constants } from "@block-feed/shared"
import { and, eq, inArray } from "drizzle-orm"
import * as schema from "@block-feed/drizzle"
import { z } from "zod"

export const zInput = z.object({
  ids: z
    .array(z.string().uuid())
    .min(constants.webhooks.limits.MAX_UUIDS.MIN)
    .max(constants.webhooks.limits.MAX_UUIDS.MAX),
})

export const handler = async (
  args: z.infer<typeof zInput>,
  ctx: GraphQLAuthContext,
) => {
  if (args.ids.length === 0) {
    return { count: 0 }
  }

  return await ctx.vendor.db.drizzle
    .delete(schema.webhook)
    .where(
      and(
        eq(schema.webhook.customerId, ctx.clerk.user.sessionClaims.sub),
        inArray(schema.webhook.id, args.ids),
      ),
    )
    .then(([result]) => ({
      count: result.affectedRows,
    }))
}
