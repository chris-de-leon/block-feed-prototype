import { type InferSelectModel } from "drizzle-orm"
import { TDatabaseLike } from "../../types"
import { webhook } from "../../schema"
import { sql } from "drizzle-orm"

export type FindOneInput = Readonly<{
  where: Readonly<Pick<InferSelectModel<typeof webhook>, "id" | "customerId">>
}>

export const findOne = async (db: TDatabaseLike, args: FindOneInput) => {
  const inputs = {
    placeholders: {
      id: sql.placeholder(webhook.id.name).getSQL(),
      customerId: sql.placeholder(webhook.customerId.name).getSQL(),
    },
    values: {
      [webhook.id.name]: args.where.id,
      [webhook.customerId.name]: args.where.customerId,
    },
  }

  return await db.query.webhook
    .findFirst({
      where(fields, operators) {
        return operators.and(
          operators.eq(fields.id, inputs.placeholders.id),
          operators.eq(fields.customerId, inputs.placeholders.customerId),
        )
      },
    })
    .prepare("webhook:find-one")
    .execute(inputs.values)
}
