# Raydium Parser Environment Setup

Your Raydium Parser environment is now set up successfully! ✅

## Environment Variables

The following environment variables are configured:
- `SOLANA_WALLET_PATH`: `/home/fidel-wole/Desktop/workspace/projects/raydium-parser/test-wallets/test-wallet.json`
- `SOLANA_RPC_ENDPOINT`: `https://api.mainnet-beta.solana.com`

## Available Commands

### Basic Usage
```bash
# Show help
./raydium-parser-env help

# Run instruction builder tests (offline mode)
./raydium-parser-env test

# Parse real transactions from Solana mainnet
./raydium-parser-env

# Run with Go
go run . test
go run . help
go run .
```

### Testing
```bash
# Run all tests
SOLANA_WALLET_PATH="$(pwd)/test-wallets/test-wallet.json" SOLANA_RPC_ENDPOINT="https://api.mainnet-beta.solana.com" go test -v

# Run specific tests
go test -v -run TestSwapInstructionBuilder
go test -v -run TestBuyInstructionBuilder
go test -v -run TestSellInstructionBuilder
go test -v -run TestCreateTokenInstructionBuilder
go test -v -run TestMigrateInstructionBuilder
go test -v -run TestBuilderChaining
```

### Load Environment Variables
```bash
# Load from .env file
source env-source.sh

# Or export manually
export SOLANA_WALLET_PATH="/home/fidel-wole/Desktop/workspace/projects/raydium-parser/test-wallets/test-wallet.json"
export SOLANA_RPC_ENDPOINT="https://api.mainnet-beta.solana.com"
```

## Test Results

✅ **Working Tests:**
- Swap Instruction Builder
- Buy Instruction Builder  
- Sell Instruction Builder
- Create Token Instruction Builder
- Migrate Instruction Builder
- Builder Method Chaining
- Basic Transaction Parsing

⚠️ **Tests with Issues:**
- Transaction Submission (RPC endpoint compatibility)
- Live Transaction Parsing (version mismatch)
- Sample Data Parsing (invalid base64)

## Next Steps

1. **For Transaction Submission Testing:**
   - Fund the test wallet with SOL
   - Use a compatible RPC endpoint (like Helius, Alchemy, or QuickNode)
   - Test with devnet first

2. **For Live Transaction Parsing:**
   - Use RPC endpoints with `maxSupportedTransactionVersion: 0`
   - Update the RPC client configuration

3. **For Sample Data:**
   - Add valid base64 encoded transaction samples to `sample_transaction.txt`

## Project Structure

```
raydium-parser/
├── main.go              # Main application
├── instructions.go      # Instruction builders
├── instructions_test.go # Test suite
├── parser.go           # Transaction parser
├── types.go            # Type definitions
├── utils.go            # Utilities
├── test-wallets/       # Test wallets
├── .env               # Environment variables
├── env-source.sh      # Environment loader
└── setup-env.sh       # Setup script
```

The environment is ready for development! 🚀
