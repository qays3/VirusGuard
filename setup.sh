#!/bin/bash

BLUE='\033[0;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
BOLD='\033[1m'
NC='\033[0m'  

function showProgressBar {
    echo -n -e "${BOLD}${BLUE}Installing dependencies... [${NC}"
    for i in {1..50}; do
        sleep 0.1
        echo -n "="
    done
    echo -e "] ${BOLD}${GREEN}Done.${NC}"
}

sudo apt update
showProgressBar

sudo apt install -y docker.io
sudo systemctl start docker
sudo systemctl enable docker
showProgressBar

sudo apt-get install -y jq
showProgressBar

chmod +x ./process/docker_containment.sh
chmod +x ./process/signature_control.sh
chmod +x ./process/terminate_process.sh

echo -e "${BOLD}${GREEN}Setup complete. Docker and dependencies installed successfully.${NC}"
