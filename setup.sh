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
sudo apt-get install jq

echo "Docker has been installed and started."
