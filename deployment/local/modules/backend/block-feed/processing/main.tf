terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

resource "docker_container" "webhook_redis" {
  count   = var.shards
  name    = "webhook-redis-${var.chain_id}-${count.index}"
  image   = var.redis_image
  command = ["--port", "6379", "--loglevel", "debug"]
  ports {
    internal = 6379
    external = var.port + count.index
  }
  networks_advanced {
    name = var.network_name
  }
}

resource "null_resource" "insert_nodes" {
  count = var.shards

  triggers = {
    // Ensures that the local-exec command is only re-run when the number of shards changes
    // NOTE: decreasing the number of shards will result in the deleted nodes still remaining in the database
    always_run = var.shards
  }

  provisioner "local-exec" {
    command = "docker exec mysql /bin/bash -c 'bash /db/utils/insert-node.sh ${docker_container.webhook_redis[count.index].name}:${docker_container.webhook_redis[count.index].ports[0].internal} ${var.chain_id}'"
  }
}

resource "docker_image" "webhook_flusher" {
  name         = "webhook-flusher:${var.tag}"
  keep_locally = true
  build {
    context    = "${path.cwd}/apps/workers"
    dockerfile = "./Dockerfile"
    build_args = {
      BUILD_PATH = "./src/apps/processing/webhook-flusher/main.go"
    }
  }
}

resource "docker_container" "webhook_flusher" {
  count   = var.shards
  name    = "webhook-flusher-${var.chain_id}-${count.index}"
  image   = docker_image.webhook_flusher.name
  restart = "always"
  env = [
    "WEBHOOK_FLUSHER_REDIS_WEBHOOK_STREAM_URL=${docker_container.webhook_redis[count.index].name}:${docker_container.webhook_redis[count.index].ports[0].internal}",
    "WEBHOOK_FLUSHER_REDIS_BLOCK_STREAM_URL=${var.etl_redis_url}",
    "WEBHOOK_FLUSHER_BLOCKCHAIN_ID=${var.chain_id}"
  ]
  networks_advanced {
    name = var.network_name
  }
}

resource "docker_image" "webhook_activator" {
  name         = "webhook-activator:${var.tag}"
  keep_locally = true
  build {
    context    = "${path.cwd}/apps/workers"
    dockerfile = "./Dockerfile"
    build_args = {
      BUILD_PATH = "./src/apps/processing/webhook-activator/main.go"
    }
  }
}

resource "docker_container" "webhook_activator" {
  count   = var.shards * var.activators_per_shard
  name    = "webhook-activator-${var.chain_id}-${count.index}"
  image   = docker_image.webhook_activator.name
  restart = "always"
  env = [
    "WEBHOOK_ACTIVATOR_MYSQL_URL=${var.mysql_url}",
    "WEBHOOK_ACTIVATOR_REDIS_URL=${docker_container.webhook_redis[count.index % var.shards].name}:${docker_container.webhook_redis[count.index % var.shards].ports[0].internal}",
    "WEBHOOK_ACTIVATOR_MYSQL_CONN_POOL_SIZE=${var.mysql_activator_conn_pool_size}",
    "WEBHOOK_ACTIVATOR_POOL_SIZE=${var.consumers_per_activator}",
    "WEBHOOK_ACTIVATOR_NAME=webhook-activator-replica-${count.index}"
  ]
  networks_advanced {
    name = var.network_name
  }
}

resource "docker_image" "webhook_consumer" {
  name         = "webhook-consumer:${var.tag}"
  keep_locally = true
  build {
    context    = "${path.cwd}/apps/workers"
    dockerfile = "./Dockerfile"
    build_args = {
      BUILD_PATH = "./src/apps/processing/webhook-consumer/main.go"
    }
  }
}

resource "docker_container" "webhook_consumer" {
  count   = var.shards * var.processors_per_shard
  name    = "webhook-consumer-${var.chain_id}-${count.index}"
  image   = docker_image.webhook_consumer.name
  restart = "always"
  env = [
    "WEBHOOK_CONSUMER_POSTGRES_URL=${var.timescaledb_url}",
    "WEBHOOK_CONSUMER_MYSQL_URL=${var.mysql_url}",
    "WEBHOOK_CONSUMER_REDIS_URL=${docker_container.webhook_redis[count.index % var.shards].name}:${docker_container.webhook_redis[count.index % var.shards].ports[0].internal}",
    "WEBHOOK_CONSUMER_MYSQL_CONN_POOL_SIZE=${var.mysql_processor_conn_pool_size}",
    "WEBHOOK_CONSUMER_POOL_SIZE=${var.consumers_per_processor}",
    "WEBHOOK_CONSUMER_NAME=webhook-consumer-replica-${count.index}"
  ]
  networks_advanced {
    name = var.network_name
  }
}

