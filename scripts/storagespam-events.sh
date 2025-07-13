#!/bin/bash

# Utility script to analyze RandomForGas events from StorageSpam contract
# Usage: ./storagespam-events.sh <RPC_URL> <CONTRACT_ADDRESS> [BATCH_SIZE]

if [ $# -lt 2 ]; then
    echo "Usage: $0 <RPC_URL> <CONTRACT_ADDRESS> [BATCH_SIZE]"
    echo ""
    echo "Example:"
    echo "  $0 https://rpc.perf-devnet-2.ethpandaops.io/ 0xFa3CE7108b73FA44a798A3aa23523c974ed5a6dE"
    echo ""
    echo "Arguments:"
    echo "  RPC_URL          - Ethereum RPC endpoint URL"
    echo "  CONTRACT_ADDRESS - StorageSpam contract address"
    echo "  BATCH_SIZE       - Number of blocks to query at once (default: 100)"
    exit 1
fi

RPC_URL="$1"
CONTRACT_ADDRESS="$2"
BATCH_SIZE="${3:-100}"

# Change to script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR" || exit 1

# Run the Go script
go run storagespam-events/storagespam-events.go -rpc "$RPC_URL" -contract "$CONTRACT_ADDRESS" -batch "$BATCH_SIZE"