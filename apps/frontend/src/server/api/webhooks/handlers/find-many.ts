import { WebhookStatus } from "@block-feed/shared/enums/webhook-status.enum"
import { AuthContext } from "@block-feed/server/graphql/types"
import * as schema from "@block-feed/drizzle"
import { z } from "zod"
import {
  InferSelectModel,
  like,
  desc,
  asc,
  and,
  eq,
  or,
  gt,
  lt,
} from "drizzle-orm"

export const zInput = z.object({
  filters: z.object({
    and: z
      .object({
        blockchain: z
          .object({
            eq: z.string().optional().nullable(),
          })
          .nullable()
          .optional(),
        status: z
          .object({
            eq: z.nativeEnum(WebhookStatus).optional().nullable(),
          })
          .nullable()
          .optional(),
        url: z
          .object({
            like: z.string().optional().nullable(),
          })
          .nullable()
          .optional(),
      })
      .nullable()
      .optional(),
  }),
  pagination: z.object({
    limit: z.number().int(),
    cursor: z
      .object({
        id: z.string().uuid(),
        reverse: z.boolean(),
      })
      .nullable()
      .optional(),
  }),
})

export const handler = async (
  args: z.infer<typeof zInput>,
  ctx: AuthContext,
) => {
  let cursor:
    | {
        webhook: InferSelectModel<typeof schema.webhook>
        data: NonNullable<(typeof args)["pagination"]["cursor"]>
      }
    | undefined = undefined

  if (args.pagination.cursor != null) {
    const webhook = await ctx.db.drizzle.query.webhook.findFirst({
      where: and(
        eq(schema.webhook.customerId, ctx.user.sub),
        eq(schema.webhook.id, args.pagination.cursor.id),
      ),
    })
    if (webhook != null) {
      cursor = {
        data: args.pagination.cursor,
        webhook,
      }
    }
  }

  return await ctx.db.drizzle.query.webhook
    .findMany({
      where: and(
        eq(schema.webhook.customerId, ctx.user.sub),
        args.filters.and?.blockchain?.eq != null &&
          args.filters.and.blockchain.eq !== ""
          ? eq(schema.webhook.blockchainId, args.filters.and.blockchain.eq)
          : undefined,
        args.filters.and?.url?.like != null && args.filters.and.url.like !== ""
          ? like(schema.webhook.url, `%${args.filters.and.url.like}%`)
          : undefined,
        args.filters.and?.status?.eq != null
          ? args.filters.and.status.eq === WebhookStatus.INACTIVE
            ? and(
                eq(schema.webhook.isActive, 0),
                eq(schema.webhook.isQueued, 0),
              )
            : args.filters.and.status.eq === WebhookStatus.PENDING
              ? and(
                  eq(schema.webhook.isActive, 0),
                  eq(schema.webhook.isQueued, 1),
                )
              : args.filters.and.status.eq === WebhookStatus.ACTIVE
                ? and(
                    eq(schema.webhook.isActive, 1),
                    eq(schema.webhook.isQueued, 1),
                  )
                : undefined
          : undefined,
        cursor != null
          ? cursor.data.reverse
            ? or(
                gt(schema.webhook.createdAt, cursor.webhook.createdAt),
                and(
                  eq(schema.webhook.createdAt, cursor.webhook.createdAt),
                  gt(schema.webhook.id, cursor.webhook.id),
                ),
              )
            : or(
                lt(schema.webhook.createdAt, cursor.webhook.createdAt),
                and(
                  eq(schema.webhook.createdAt, cursor.webhook.createdAt),
                  lt(schema.webhook.id, cursor.webhook.id),
                ),
              )
          : undefined,
      ),
      limit: args.pagination.limit + 1, // the +1 helps us determine if we have more data for pagination
      orderBy: cursor?.data.reverse
        ? [asc(schema.webhook.createdAt), asc(schema.webhook.id)]
        : [desc(schema.webhook.createdAt), desc(schema.webhook.id)],
    })
    .then((result) => {
      const reverse = cursor?.data.reverse ?? false
      const isFirstPage = cursor == null

      if (result.length > args.pagination.limit) {
        result.pop()
        return {
          payload: reverse ? result.reverse() : result,
          pagination: {
            hasNext: true,
            hasPrev: !isFirstPage,
          },
        }
      }

      return {
        payload: reverse ? result.reverse() : result,
        pagination: {
          hasNext: isFirstPage ? false : reverse,
          hasPrev: isFirstPage ? false : !reverse,
        },
      }
    })
}
