set -e

. ./tools/utils/utils.sh

export_env_files ./env/dev
echo ""

find ./libs/block-gateway -name '*.tests.ts' | tr '\n' ' ' | xargs node --require ts-node/register --test --test-reporter spec
