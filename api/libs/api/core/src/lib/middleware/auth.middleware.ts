import { database } from "@api/shared/database"
import { MiddlewareOpts } from "../types"
import { auth0 } from "@api/shared/auth0"
import { TRPCError } from "@trpc/server"
import { sql } from "drizzle-orm"

export const requireAuth = async (
  opts: MiddlewareOpts<
    Readonly<{
      database: ReturnType<typeof database.core.createClient>
      auth0: ReturnType<typeof auth0.createClient>
    }>
  >,
) => {
  // Gets the authorization header value
  const authorization =
    opts.ctx.event.headers["authorization"] ??
    opts.ctx.event.headers["Authorization"]

  // Checks that the authorization header exists
  if (authorization == null) {
    throw new TRPCError({
      code: "BAD_REQUEST",
      message: "request is missing authorization header",
    })
  }

  // Checks that the authorization header value starts with 'bearer' (case insensitive)
  const value = authorization.trim()
  if (!value.toLowerCase().startsWith("bearer")) {
    throw new TRPCError({
      code: "BAD_REQUEST",
      message: 'authorization header value is missing "bearer" prefix',
    })
  }

  // Parses the authorization header
  const tokens = value.split(" ")
  if (tokens.length <= 0) {
    throw new TRPCError({
      code: "BAD_REQUEST",
      message: "authorization header value is malformed",
    })
  }

  // Extracts the auth token
  const accessToken = tokens[tokens.length - 1]

  // Uses the auth token to get the user profile info
  // TODO: cache info somewhere?
  const profile = await opts.ctx.auth0.userInfo
    .getUserInfo(accessToken)
    .then(({ data }) => data)
    .catch(() => {
      throw new TRPCError({
        code: "UNAUTHORIZED",
        message: "invalid access token",
      })
    })

  // Prepares insert parameters
  const inputs = {
    placeholders: {
      id: sql.placeholder(database.schema.customer.id.name),
    },
    values: {
      [database.schema.customer.id.name]: profile.sub,
    },
  }

  // Inserts the user (or ignores if one already exists)
  await opts.ctx.database.drizzle
    .insert(database.schema.customer)
    .values({ id: inputs.placeholders.id })
    .onConflictDoNothing({
      target: database.schema.customer.id,
    })
    .prepare("customer:upsert")
    .execute(inputs.values)

  // Adds the auth0 profile info to the context
  return opts.next({
    ctx: { user: profile },
  })
}
