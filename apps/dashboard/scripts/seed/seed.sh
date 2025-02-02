#!/bin/bash
set -e

if [ -z "$1" ]; then
  echo "argument 1 is required (User ID)"
  exit 1
fi

if [ -z "$2" ]; then
  echo "argument 2 is required (start port)"
  exit 1
fi

if [ -z "$3" ]; then
  echo "argument 3 is required (end port)"
  exit 1
fi

# Inserts the user
docker exec -it mysql-dev /bin/bash -c '
  mysql --password="password" --database="dev" -e"INSERT IGNORE INTO customer VALUES (\"'$1'\", DEFAULT)"
'

# Inserts blockchains and redis cluster info
chains=("flow" "ethereum")
port=7001
for ((i = 0; i < ${#chains[@]}; i++)); do
  cport=$((port + (1000 * i)))
  chain=${chains[$i]}
  url=${chain_urls[$i]}
  bash ./scripts/seed/blockchains/insert.sh \
    "$chain" \
    "1283" \
    "http://dummy-chain-url:3000" \
    "http://dummy-pg-store-url:5432" \
    "http://dummy-redis-store-url:6379" \
    "localhost:$cport" \
    "http://dummy-redis-stream-url:6379"
done
bash ./scripts/seed/blockchains/list.sh

# Inserts some webhooks
num_chains=${#chains[@]}
for ((i = $2; i <= $3; i++)); do
  rand_idx=$(($RANDOM % $num_chains))
  rand_chain=${chains[$rand_idx]}
  bash ./scripts/seed/webhooks/insert.sh \
    "http://localhost:$i" \
    "$1" \
    "$rand_chain"
done
bash ./scripts/seed/webhooks/list.sh
