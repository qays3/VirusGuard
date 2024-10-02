#!/bin/bash

MALWARE_NAME=$1

if [ -z "$MALWARE_NAME" ]; then
    echo "No malware name provided."
    exit 1
fi

PIDS=$(pgrep -f "$MALWARE_NAME")

if [ -z "$PIDS" ]; then
    echo "No processes found for $MALWARE_NAME."
else
    echo "Terminating processes for $MALWARE_NAME with PIDs: $PIDS"
    kill $PIDS
    echo "Processes terminated."
fi