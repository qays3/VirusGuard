#!/bin/bash

MALWARE_NAME=$1
CONTAINER_NAME=$2
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'  
BOLD='\033[1m'
NC='\033[0m'  

if [ -z "$MALWARE_NAME" ] || [ -z "$CONTAINER_NAME" ]; then
    echo "Usage: $0 <malware_name> <container_name>"
    exit 1
fi

docker rm -f "$CONTAINER_NAME" 2>/dev/null

echo "Starting Docker container: $CONTAINER_NAME"

docker run --name "$CONTAINER_NAME" -d -it --rm ubuntu:latest /bin/bash -c "mkdir -p /malware && /bin/bash"
if [ $? -ne 0 ]; then
    echo "Failed to run Docker container"
    exit 1
fi

echo "Copying malware to container"
docker cp "$MALWARE_NAME" "$CONTAINER_NAME:/malware/$MALWARE_NAME"
if [ $? -ne 0 ]; then
    echo "Failed to copy malware to container"
    exit 1
fi

echo "Executing malware in container"
docker exec "$CONTAINER_NAME" /malware/"$MALWARE_NAME"
if [ $? -ne 0 ]; then
    echo "Failed to execute malware in container"
    exit 1
fi

echo -e "Checking container status:\n"
docker ps -a | grep "$CONTAINER_NAME"

 
echo -e "\n${GREEN}${BOLD}Run the following command to enter the container:${NC}"
echo -e "${BLUE}${BOLD}docker exec -it $CONTAINER_NAME /bin/bash${NC}"  
echo -e "${GREEN}${BOLD}You will find the malware inside the /malware folder.${NC}"
