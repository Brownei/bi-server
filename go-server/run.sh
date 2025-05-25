#!/bin/bash

set -e

CONTAINER_NAME="tcp_db"
IMAGE_NAME="postgres:alpine"
POSTGRES_USER="tcp_db"
POSTGRES_PASSWORD="tcp_db"
POSTGRES_DB="tcp_db"

# Check if container exists
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    # Container exists
    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        echo "Container '$CONTAINER_NAME' is already running."
    else
        echo "Container '$CONTAINER_NAME' exists but is not running. Restarting..."
        docker restart "$CONTAINER_NAME"
    fi
else
    echo "Container '$CONTAINER_NAME' does not exist. Creating it..."
    docker run -d \
        --name "$CONTAINER_NAME" \
        -p 5432:5432 \
        -e POSTGRES_USER="$POSTGRES_USER" \
        -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
        -e POSTGRES_DB="$POSTGRES_DB" \
        "$IMAGE_NAME"
fi

# Wait for PostgreSQL to be ready
echo "Waiting for database instance to be ready..."
until docker exec "$CONTAINER_NAME" pg_isready -U "$POSTGRES_USER" > /dev/null 2>&1; do
    sleep 1
done
echo "Database instance is ready."

echo "Running project..."
make run

