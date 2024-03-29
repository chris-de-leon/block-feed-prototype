set -e

# Sample usages:
#
#   pnpm run api:test
#
#   pnpm run api:test -- refresh
#

. ./scripts/utils/utils.sh

export_env_files .
echo ""

if [ ! -z $1 ]; then
	printf "\n\nGenerating Drizzle schema...\n\n"
	pnpm run db:introspect

	printf "\n\nGenerating GraphQL SDK...\n\n"
	pnpm run codegen
fi

printf "\n\nRunning tests...\n\n"
find ./tests -type f -name '*.tests.ts' |
	tr '\n' ' ' |
	xargs node \
		--require ts-node/register \
		--require tsconfig-paths/register \
		--test \
		--test-concurrency 1 \
		--test-reporter \
		spec
