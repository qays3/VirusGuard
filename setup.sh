#!/bin/bash

function showProgressBar {
    echo -n "Installing dependencies... ["
    for i in {1..50}; do
        sleep 0.1
        echo -n "="
    done
    echo "] Done."
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

echo "Setup complete. Docker and dependencies installed successfully."
