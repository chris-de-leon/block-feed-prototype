set -e

. ./tools/utils/utils.sh
. ./tools/db/utils.sh

# Define an array of required options and their associated variable names
declare -A required_options=(
  ["--environment"]="environment"
)

# Parse command-line options
while [ $# -gt 0 ]; do
  env_var_name="${required_options[$1]}"
  if [ -z "$env_var_name" ]; then
    echo "unknown option: $1"
    exit 1
  fi

  eval "$env_var_name=\"$2\""
  shift 2
done

# Checks that required flags were provided
for flag in "${!required_options[@]}"; do
  env_var_name="${required_options[$flag]}"
  env_var_val="${!env_var_name}"
  if [ -z "$env_var_val" ]; then
    echo "error: $flag is required"
    exit 1
  fi
done

# Run migrations
if [ "$environment" == "stag" ] || [ "$environment" == "prod" ]; then
  # We're using an SSH tunnel to connect from localhost to
  # an EC2 instance to the RDS instance. The communication
  # from localhost to EC2 is protected by default since we
  # are using SSH. The communication from EC2 to RDS is also
  # secured by security groups and firewalls. As a result,
  # we don't need to provide SSL certs to connect to RDS.
  #
  # https://stackoverflow.com/a/45088585/22520608
  #
  # TODO: test migrations on planetscale
  #
  export NODE_TLS_REJECT_UNAUTHORIZED='0'
  export DRIZZLE_DB_MIGRATIONS_FOLDER="./drizzle/migrations/$environment"
  export DRIZZLE_DB_MODE="planetscale"

  export_env_files "./env/$environment"
  create_temp_db_tunnel "$environment" "5430"
  drizzle-kit generate:mysql
  ts-node ./drizzle/scripts/migrate.ts
  ts-node ./drizzle/scripts/refresh-roles.ts
else
  export DRIZZLE_DB_MIGRATIONS_FOLDER="./drizzle/migrations/dev"
  export DRIZZLE_DB_MODE="default"

  export DRIZZLE_DB_NAME="dev" # `block_feed` database may not exist yet, so we need to use this one
  export_env_files "./env/dev"
  ts-node ./drizzle/scripts/recreate-database.ts

  export DRIZZLE_DB_NAME="block_feed"
  drizzle-kit push:mysql
  ts-node ./drizzle/scripts/refresh-roles.ts
fi
