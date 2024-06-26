import { defineConfig } from "drizzle-kit"

export default defineConfig({
  out: "./src/generated",
  dialect: "mysql",
  verbose: true,
  strict: true,
  dbCredentials: {
    url: process.env.DRIZZLE_DB_URL ?? "",
  },
  introspect: {
    casing: "camel",
  },
})
