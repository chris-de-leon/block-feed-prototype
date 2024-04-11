import { gqlUnauthorizedError, gqlBadRequestError } from "../errors"
import { db } from "@block-feed/server/vendor/database"
import { auth } from "@block-feed/server/vendor/auth0"
import { ApiCache } from "@block-feed/server/caching"
import * as schema from "@block-feed/drizzle"
import { UserInfoResponse } from "auth0"

export type RequireAuthContext = Readonly<{
  cache: ApiCache<UserInfoResponse>
  auth0: ReturnType<typeof auth.client.create>
  db: ReturnType<typeof db.client.create>
  req: Request
}>

export const requireAuth = async (ctx: RequireAuthContext) => {
  // Gets the authorization header value
  const authorization =
    ctx.req.headers.get("authorization") ?? ctx.req.headers.get("Authorization")

  // Checks that the authorization header exists
  if (authorization == null) {
    throw gqlBadRequestError("request is missing authorization header")
  }

  // Checks that the authorization header is a string
  if (typeof authorization !== "string") {
    throw gqlBadRequestError("authorization header cannot have multiple values")
  }

  // Checks that the authorization header value starts with 'bearer' (case insensitive)
  const value = authorization.trim()
  if (!value.toLowerCase().startsWith("bearer")) {
    throw gqlBadRequestError(
      'authorization header value is missing "bearer" prefix',
    )
  }

  // Parses the authorization header
  const tokens = value.split(" ")
  if (tokens.length <= 0) {
    throw gqlBadRequestError("authorization header value is malformed")
  }

  // Extracts the auth token
  const accessToken = tokens[tokens.length - 1]

  // Check the cache for the user's profile info
  const cached = await ctx.cache.get(accessToken)
  if (cached != null) {
    return cached
  }

  // Uses the auth token to get the user profile info
  const profile = await ctx.auth0.userInfo
    .getUserInfo(accessToken)
    .then(({ data }) => data)
    .catch(() => null)

  // If the access token is not valid, return an error
  // Otherwise cache the info for next time
  if (profile == null) {
    throw gqlUnauthorizedError("invalid access token")
  } else {
    await ctx.cache.set(accessToken, profile)
  }

  // Inserts the user in the database (or ignores if one already exists)
  await ctx.db.drizzle
    .insert(schema.customer)
    .ignore()
    .values({ id: profile.sub })

  // Adds the auth0 profile info to the context
  return profile
}