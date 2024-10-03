#!/bin/bash

if [ "$#" -lt 3 ]; then
    echo "Usage: $0 <malware_name> <signature> <action>"
    echo "Actions: block | unblock"
    exit 1
fi

MALWARE_NAME=$1
SIGNATURE=$2
ACTION=$3
SIGNATURE_FILE="./json/signature.json"

mkdir -p "$(dirname "$SIGNATURE_FILE")"

if [ ! -f "$SIGNATURE_FILE" ]; then
    echo "[]" > "$SIGNATURE_FILE"
fi

block_signature() {
    if grep -q "\"malware\":\"$MALWARE_NAME\"" "$SIGNATURE_FILE"; then
        echo "Malware $MALWARE_NAME is already blocked."
    else
        jq --arg name "$MALWARE_NAME" --arg sig "$SIGNATURE" '. += [{"malware": $name, "signature": $sig}]' "$SIGNATURE_FILE" > tmp.$$.json && mv tmp.$$.json "$SIGNATURE_FILE"
        echo "Malware $MALWARE_NAME has been blocked with signature: $SIGNATURE."
    fi
}

unblock_signature() {
    if ! grep -q "\"malware\":\"$MALWARE_NAME\"" "$SIGNATURE_FILE"; then
        echo "Error: $MALWARE_NAME is not blocked."
        exit 1
    fi

    jq --arg name "$MALWARE_NAME" 'map(select(.malware != $name))' "$SIGNATURE_FILE" > tmp.$$.json && mv tmp.$$.json "$SIGNATURE_FILE"
    echo "Malware $MALWARE_NAME has been unblocked."
}

case $ACTION in
    block)
        block_signature
        ;;
    unblock)
        unblock_signature
        ;;
    *)
        echo "Invalid action: $ACTION. Use 'block' or 'unblock'."
        exit 1
        ;;
esac
