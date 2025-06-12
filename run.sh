  #!/bin/bash

set -e

run() {
# Postgres container configuration
CONTAINER_NAME="tcp_db"
IMAGE_NAME="postgres:alpine"
POSTGRES_USER="tcp_db"
POSTGRES_PASSWORD="tcp_db"
POSTGRES_DB="tcp_db"

# Paths
GO_SERVER_DIR="./go-server"
RUST_CLIENT_DIR="./rust-client"
NETWORK_NAME="tcp_network"

# 0. Create a network firstly before anything
if ! docker network ls --format '{{.Name}}' | grep -q "^${NETWORK_NAME}$"; then
    echo "Creating Docker network '$NETWORK_NAME'..."
    docker network create "$NETWORK_NAME"
fi

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
        --network "$NETWORK_NAME" \
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
echo "Checking for existing 'go-server' Docker image..."

if docker images --format '{{.Repository}}:{{.Tag}}' | grep -q '^go-server:latest$'; then
    echo "Finding containers using the 'go-server' image..."
    container_ids=$(docker ps -a --filter "ancestor=go-server" --format '{{.ID}}')

    if [ -n "$container_ids" ]; then
        echo "Stopping and removing containers using 'go-server'..."
        for id in $container_ids; do
            docker stop "$id"
            docker rm "$id"
        done
    fi

    echo "Removing 'go-server' Docker image..."
    docker rmi -f go-server
fi


echo "Building new 'go-server' Docker image..."
docker build --no-cache -t go-server "$GO_SERVER_DIR"

if docker ps -a --format '{{.Names}}' | grep -q '^go_server$'; then
    echo "Stopping and removing existing 'go_server' container..."
    docker stop go_server
    docker rm go_server
fi

echo "Running new 'go_server' container..."
docker run --rm -d --name go_server --network "$NETWORK_NAME" -p 3000:3000 go-server

# 4. Run Rust client (local cargo run)
echo "Running Rust client with cargo..."
pushd "$RUST_CLIENT_DIR" > /dev/null
cargo run
popd > /dev/null
}

clean() {
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
}

case "$1" in 
  run)
    run 
    exit 0
    ;;
  clean)
    clean 
    exit 0
    ;;
esac

echo "Invalid command: $1"
