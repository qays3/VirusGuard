#!/bin/bash

MALWARE_NAME=$1
CONTAINER_NAME=$2

if [ -z "$MALWARE_NAME" ] || [ -z "$CONTAINER_NAME" ]; then
    echo "Usage: $0 <malware_name> <container_name>"
    exit 1
fi

mkdir -p "$CONTAINER_NAME"
docker run --name "$CONTAINER_NAME" -d -it --rm -v "$(pwd)/$CONTAINER_NAME:/malware" ubuntu:latest /bin/bash

docker cp "$MALWARE_NAME" "$CONTAINER_NAME:/malware/$MALWARE_NAME"
docker exec "$CONTAINER_NAME" /malware/"$MALWARE_NAME"