#!/bin/bash

# Source this file to set up environment variables
# Usage: source env-source.sh

if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | grep -v '^$' | xargs)
    echo "Environment variables loaded from .env file"
    echo "SOLANA_WALLET_PATH: $SOLANA_WALLET_PATH"
    echo "SOLANA_RPC_ENDPOINT: $SOLANA_RPC_ENDPOINT"
else
    echo "Error: .env file not found"
    return 1
fi

# Check if wallet file exists
if [ -f "$SOLANA_WALLET_PATH" ]; then
    echo "✓ Test wallet found"
else
    echo "✗ Test wallet not found at: $SOLANA_WALLET_PATH"
    return 1
fi

echo ""
echo "Environment ready! You can now run:"
echo "  go run . test           # Run instruction builder demos"
echo "  go run .                # Parse real transactions"
echo "  go test -v              # Run all tests"
