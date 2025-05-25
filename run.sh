#!/bin/bash

set -e

# Postgres container configuration
CONTAINER_NAME="tcp_db"
IMAGE_NAME="postgres:alpine"
POSTGRES_USER="tcp_db"
POSTGRES_PASSWORD="tcp_db"
POSTGRES_DB="tcp_db"

# Paths
GO_SERVER_DIR="./go-server"
RUST_CLIENT_DIR="./rust-client"

# 1. Start or restart Postgres container
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        echo "Postgres container '$CONTAINER_NAME' is already running."
    else
        echo "Restarting existing Postgres container..."
        docker restart "$CONTAINER_NAME"
    fi
else
    echo "Creating new Postgres container '$CONTAINER_NAME'..."
    docker run -d \
        --name "$CONTAINER_NAME" \
        -p 5432:5432 \
        -e POSTGRES_USER="$POSTGRES_USER" \
        -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
        -e POSTGRES_DB="$POSTGRES_DB" \
        "$IMAGE_NAME"
fi

# 2. Wait for Postgres to be ready
echo "Waiting for PostgreSQL to be ready..."
until docker exec "$CONTAINER_NAME" pg_isready -U "$POSTGRES_USER" > /dev/null 2>&1; do
    sleep 1
done
echo "PostgreSQL is ready."

# 3. Build and run Go server in Docker
echo "Building and running Go server from Dockerfile..."
docker build -t go-server "$GO_SERVER_DIR"
docker run --rm -d --name go_server --network host go-server

# 4. Run Rust client (local cargo run)
echo "Running Rust client with cargo..."
pushd "$RUST_CLIENT_DIR" > /dev/null
cargo run
popd > /dev/null

