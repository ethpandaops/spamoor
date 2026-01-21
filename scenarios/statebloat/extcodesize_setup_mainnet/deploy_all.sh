#!/bin/bash

# Deploy all extcodesize contracts sequentially
# Usage: ./deploy_all.sh [RPC_URL] [PRIVATE_KEY]

RPC_URL="${1:-http://localhost:8545}"
PRIVATE_KEY="${2:-0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SPAMOOR_BIN="${SCRIPT_DIR}/../../../bin/spamoor"

echo "=== EXTCODESIZE Contract Deployment ==="
echo "RPC: $RPC_URL"
echo ""

# # Deploy 0.5KB contracts
# echo "[1/8] Deploying 0.5KB contracts..."
# $SPAMOOR_BIN run -h "$RPC_URL" -p "$PRIVATE_KEY" "$SCRIPT_DIR/deploy_0_5kb.yaml"
# echo ""

# # Deploy 1KB contracts
# echo "[2/8] Deploying 1KB contracts..."
# $SPAMOOR_BIN run -h "$RPC_URL" -p "$PRIVATE_KEY" "$SCRIPT_DIR/deploy_1kb.yaml"
# echo ""

# # Deploy 2KB contracts
# echo "[3/8] Deploying 2KB contracts..."
# $SPAMOOR_BIN run -h "$RPC_URL" -p "$PRIVATE_KEY" "$SCRIPT_DIR/deploy_2kb.yaml"
# echo ""

# Deploy 5KB contracts
echo "[4/8] Deploying 5KB contracts..."
$SPAMOOR_BIN run -h "$RPC_URL" -p "$PRIVATE_KEY" "$SCRIPT_DIR/deploy_5kb.yaml"
echo ""

# Deploy 10KB contracts
echo "[5/8] Deploying 10KB contracts..."
$SPAMOOR_BIN run -h "$RPC_URL" -p "$PRIVATE_KEY" "$SCRIPT_DIR/deploy_10kb.yaml"
echo ""

# Deploy 24KB contracts
echo "[6/8] Deploying 24KB contracts..."
$SPAMOOR_BIN run -h "$RPC_URL" -p "$PRIVATE_KEY" "$SCRIPT_DIR/deploy_24kb.yaml"
echo ""

# Deploy 32KB contracts
echo "[7/8] Deploying 32KB contracts..."
$SPAMOOR_BIN run -h "$RPC_URL" -p "$PRIVATE_KEY" "$SCRIPT_DIR/deploy_32kb.yaml"
echo ""

# Deploy 64KB contracts
echo "[8/8] Deploying 64KB contracts..."
$SPAMOOR_BIN run -h "$RPC_URL" -p "$PRIVATE_KEY" "$SCRIPT_DIR/deploy_64kb.yaml"
echo ""

echo "=== All deployments complete ==="
