terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

provider "docker" {
}

resource "docker_image" "redis" {
  name         = "docker.io/redis:${var.redis_version}"
  keep_locally = true
}

resource "docker_image" "timescaledb_dev" {
  name         = "timescaledb-dev:${var.tag}"
  keep_locally = true
  build {
    context    = "${path.cwd}/vendor/timescaledb"
    dockerfile = "./Dockerfile"
    build_args = {
      PG_VERSION = var.timescaledb_version
    }
  }
}

resource "docker_network" "block_feed_net" {
  name   = "block_feed_net"
  driver = "bridge"
}


module "block_feed_moonbeam_mainnet" {
  source                   = "./modules/backend/block-feed"
  tag                      = var.tag
  network_name             = docker_network.block_feed_net.name
  timescaledb_image        = docker_image.timescaledb_dev.name
  redis_image              = docker_image.redis.name
  chain_url                = "wss://moonbeam-rpc.dwellir.com"
  chain_id                 = "moonbeam-mainnet"
  chain_name               = "eth"
  redis_block_store_port   = 7000
  redis_cluster_start_port = 7001
  redis_cluster_end_port   = 7006
  redis_block_stream_port  = 7007
  timescaledb_port         = 7008
  shard_count              = 1
  replicas_per_shard       = 1
  workers_per_replica      = 1
}

# module "block_feed_flow_testnet" {
#   source                   = "./modules/backend/block-feed"
#   tag                      = var.tag
#   network_name             = docker_network.block_feed_net.name
#   timescaledb_image        = docker_image.timescaledb_dev.name
#   redis_image              = docker_image.redis.name
#   chain_url                = "access.devnet.nodes.onflow.org:9000"
#   chain_id                 = "flow-testnet"
#   chain_name               = "flow"
#   redis_block_store_port   = 8000
#   redis_cluster_start_port = 8001
#   redis_cluster_end_port   = 8006
#   redis_block_stream_port  = 8007
#   timescaledb_port         = 8008
#   shard_count              = 1
#   replicas_per_shard       = 1
#   workers_per_replica      = 1
# }
