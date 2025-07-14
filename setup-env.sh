#!/bin/bash

# Raydium Parser Environment Setup Script

echo "Setting up Raydium Parser environment..."

# Set environment variables
export SOLANA_WALLET_PATH="/home/fidel-wole/Desktop/workspace/projects/raydium-parser/test-wallets/test-wallet.json"
export SOLANA_RPC_ENDPOINT="https://api.mainnet-beta.solana.com"

# Check if wallet file exists
if [ -f "$SOLANA_WALLET_PATH" ]; then
    echo "✓ Test wallet found at: $SOLANA_WALLET_PATH"
else
    echo "✗ Test wallet not found at: $SOLANA_WALLET_PATH"
    exit 1
fi

# Print environment variables
echo "Environment variables set:"
echo "SOLANA_WALLET_PATH: $SOLANA_WALLET_PATH"
echo "SOLANA_RPC_ENDPOINT: $SOLANA_RPC_ENDPOINT"

# Check if Go is available
if command -v go &> /dev/null; then
    echo "✓ Go is available: $(go version)"
else
    echo "✗ Go is not available"
    exit 1
fi

# Check if project dependencies are available
echo "Checking dependencies..."
if go mod tidy; then
    echo "✓ Dependencies are up to date"
else
    echo "✗ Failed to update dependencies"
    exit 1
fi

# Try to build the project
echo "Building project..."
if go build -o raydium-parser-test; then
    echo "✓ Project built successfully"
    rm -f raydium-parser-test
else
    echo "✗ Failed to build project"
    exit 1
fi

echo ""
echo "🚀 Environment setup complete!"
echo ""
echo "Key files created:"
echo "  - .env (environment variables)"
echo "  - env-source.sh (environment loader)"
echo "  - ENVIRONMENT.md (detailed setup guide)"
echo ""
echo "Quick start:"
echo "  ./raydium-parser-env help    # Show help"
echo "  ./raydium-parser-env test    # Run instruction builder tests"
echo "  go test -v -run TestSwapInstructionBuilder  # Run specific test"
echo ""
echo "For detailed instructions, see ENVIRONMENT.md"
echo ""
echo "✅ Ready to develop with Raydium Parser!"
