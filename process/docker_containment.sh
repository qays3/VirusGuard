#!/bin/bash

MALWARE_NAME=$1
CONTAINER_NAME=$2

if [ -z "$MALWARE_NAME" ] || [ -z "$CONTAINER_NAME" ]; then
    echo "Usage: $0 <malware_name> <container_name>"
    exit 1
fi

mkdir -p "$CONTAINER_NAME"
docker rm -f "$CONTAINER_NAME" 2>/dev/null


echo "Starting Docker container: $CONTAINER_NAME"
docker run --name "$CONTAINER_NAME" -d -it --rm -v "$(pwd)/$CONTAINER_NAME:/malware" ubuntu:latest /bin/bash
if [ $? -ne 0 ]; then
    echo "Failed to run Docker container"
    exit 1
fi

docker cp "$MALWARE_NAME" "$CONTAINER_NAME:/malware/$MALWARE_NAME"
if [ $? -ne 0 ]; then
    echo "Failed to copy malware to container"
    exit 1
fi

docker exec "$CONTAINER_NAME" /malware/"$MALWARE_NAME"
if [ $? -ne 0 ]; then
    echo "Failed to execute malware in container"
    exit 1
fi
