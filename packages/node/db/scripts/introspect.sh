#!/bin/bash

# Docker configs
MYSQL_VERSION="8.3.0"
RANDM_UUID="$(uuidgen)"
IMAGE_NAME="db:$RANDM_UUID"

# MySQL configs
MYSQL_ROOT_PWORD="password"
MYSQL_UNAME="rootuser"
MYSQL_PWORD="password"
MYSQL_PORT="3306"
MYSQL_DB="drizzle"

# Defines a helper function for cleaning up resources
cleanup() {
	docker stop "$RANDM_UUID"
	docker image rm "$IMAGE_NAME"
}

# Cleans up the image, container, and temporary drizzle config file once this script exits
trap cleanup EXIT

# If any errors occur exit the script
set -e

# Builds a custom image with our setup scripts
docker build \
	-t "$IMAGE_NAME" \
	--build-arg MYSQL_VERSION="$MYSQL_VERSION" \
	"$(git rev-parse --show-toplevel)/vendor/mysql"

# Creates a database container and applies the setup scripts
docker run --rm -d \
	-e MYSQL_ROOT_PASSWORD="$MYSQL_ROOT_PWORD" \
	-e MYSQL_PASSWORD="$MYSQL_PWORD" \
	-e MYSQL_DATABASE="$MYSQL_DB" \
	-e MYSQL_USER="$MYSQL_UNAME" \
	-p 3306:3306 \
	--name "$RANDM_UUID" \
	"$IMAGE_NAME"

# Waits for the database to come online
TOTAL_SECONDS=0
PING_SECONDS=2
while ! docker exec "$RANDM_UUID" mysqladmin --user="root" --password="$MYSQL_ROOT_PWORD" --host="127.0.0.1" ping --silent &>/dev/null; do
	echo "Waiting for database connection... ($TOTAL_SECONDS seconds)"
	sleep $PING_SECONDS
	((TOTAL_SECONDS += $PING_SECONDS))
done

# Generates a drizzle schema
export DRIZZLE_DB_URL="mysql://root:$MYSQL_ROOT_PWORD@host.docker.internal:$MYSQL_PORT/$MYSQL_DB"
drizzle-kit introspect
