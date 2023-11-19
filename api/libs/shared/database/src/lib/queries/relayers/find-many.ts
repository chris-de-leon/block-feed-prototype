import { CONSTANTS, createClient } from "../../core"
import { type InferSelectModel } from "drizzle-orm"
import { TRowLimit, TRowOffset } from "../../types"
import { relayers } from "../../schema"
import { sql } from "drizzle-orm"

export type FindManyInput = Readonly<
  Pick<InferSelectModel<typeof relayers>, "userId"> & TRowLimit & TRowOffset
>

export const findMany = async (
  db: ReturnType<typeof createClient>,
  args: FindManyInput,
) => {
  const inputs = {
    placeholders: {
      userId: sql.placeholder(relayers.userId.name).getSQL(),
      limit: sql.placeholder(CONSTANTS.LIMIT).getSQL(),
      offset: sql.placeholder(CONSTANTS.OFFSET).getSQL(),
    },
    values: {
      [relayers.userId.name]: args.userId,
      [CONSTANTS.LIMIT]: args.limit,
      [CONSTANTS.OFFSET]: args.offset,
    },
  }

  return await db.drizzle.query.relayers
    .findMany({
      where(fields, operators) {
        return operators.eq(fields.userId, inputs.placeholders.userId)
      },
      orderBy(fields, operators) {
        return [operators.desc(fields.createdAt), operators.desc(fields.id)]
      },
    })
    .prepare()
    .execute(inputs.values)
}
