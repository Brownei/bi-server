#!/bin/bash

set -e

# Names
POSTGRES_CONTAINER="tcp_db"
GO_CONTAINER="go_server"
GO_IMAGE="go-server"
NETWORK_NAME="tcp_network"

echo "Stopping and removing containers..."
docker rm -f "$POSTGRES_CONTAINER" || echo "Postgres container not running."
docker rm -f "$GO_CONTAINER" || echo "Go server container not running."

echo "Removing Docker image for Go server..."
docker rmi -f "$GO_IMAGE" || echo "Go server image not found."

echo "Removing Docker network if exists..."
docker network rm "$NETWORK_NAME" || echo "Network '$NETWORK_NAME' not found."

echo "Cleanup complete."

