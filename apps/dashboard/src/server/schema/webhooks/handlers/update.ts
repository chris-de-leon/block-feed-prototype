import { constants } from "@block-feed/dashboard/utils/constants"
import { GraphQLAuthContext } from "../../../graphql/types"
import * as schema from "@block-feed/node-db"
import { and, eq } from "drizzle-orm"
import { z } from "zod"

export const zInput = z.object({
  id: z.string().uuid(),
  data: z.object({
    url: z
      .string()
      .url()
      .min(constants.webhooks.limits.URL_LEN.MIN)
      .max(constants.webhooks.limits.URL_LEN.MAX)
      .optional()
      .nullable(),
    maxBlocks: z
      .number()
      .int()
      .min(constants.webhooks.limits.MAX_BLOCKS.MIN)
      .max(constants.webhooks.limits.MAX_BLOCKS.MAX)
      .optional()
      .nullable(),
    maxRetries: z
      .number()
      .int()
      .min(constants.webhooks.limits.MAX_RETRIES.MIN)
      .max(constants.webhooks.limits.MAX_RETRIES.MAX)
      .optional()
      .nullable(),
    timeoutMs: z
      .number()
      .int()
      .min(constants.webhooks.limits.TIMEOUT_MS.MIN)
      .max(constants.webhooks.limits.TIMEOUT_MS.MAX)
      .optional()
      .nullable(),
  }),
})

export const handler = async (
  args: z.infer<typeof zInput>,
  ctx: GraphQLAuthContext,
) => {
  if (Object.values(args.data).filter((v) => v != null).length === 0) {
    return { count: 0 }
  }

  return await ctx.providers.mysql.drizzle
    .update(schema.webhook)
    .set({
      maxRetries: args.data.maxRetries ?? undefined,
      maxBlocks: args.data.maxBlocks ?? undefined,
      timeoutMs: args.data.timeoutMs ?? undefined,
      url: args.data.url ?? undefined,
    })
    .where(
      and(
        eq(schema.webhook.customerId, ctx.clerk.user.id),
        eq(schema.webhook.id, args.id),
      ),
    )
    .then(([result]) => ({
      count: result.affectedRows,
    }))
}
