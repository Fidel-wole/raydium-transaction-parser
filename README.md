# Raydium Transaction Parser and Instruction Builder

A comprehensive Go library for parsing Raydium transactions from the Solana blockchain and building Raydium instructions for submission.

## Features

### Transaction Parsing
- **Dual Format Support**: Handles both Geyser and standard Solana RPC transaction formats
- **Multi-Program Support**: Recognizes all major Raydium program IDs including V4, V5, Launchpad, and CP-Swap
- **Instruction Detection**: Identifies and parses swap, buy, sell, create, and migrate operations
- **Generic Parsing**: Handles unknown instruction discriminators with fallback parsing logic
- **Real-time Analysis**: Fetches and parses live transactions from Solana mainnet

### Instruction Building
- **Builder Pattern**: Fluent API with method chaining for easy instruction construction
- **Type Safety**: Strongly typed instruction builders with validation
- **Multiple Instruction Types**:
  - `SwapInstruction`: For token swaps through Raydium AMM
  - `BuyInstruction`: For token purchases in Raydium Launchpad
  - `SellInstruction`: For token sales in Raydium Launchpad
  - `CreateTokenInstruction`: For token creation operations
  - `MigrateInstruction`: For pool migration operations
- **Solana Integration**: Built-in serialization to valid Solana instructions

## Prerequisites

- Go 1.19 or later
- Internet connection for downloading dependencies

## Installation

1. Clone or download this project
2. Navigate to the project directory
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Command Line Interface

```bash
# Run instruction builder tests
go run . test

# Parse a real transaction from Solana mainnet
go run .

# Show help
go run . help

# Run in offline mode (same as test)
go run . offline
```

### Building Instructions

```go
package main

import (
    "fmt"
    "github.com/gagliardetto/solana-go"
)

func main() {
    // Create a swap instruction
    swapInst := NewSwapInstruction().
        SetUserSourceToken(solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")).
        SetUserDestToken(solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")).
        SetUserOwner(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH")).
        SetAmountIn(1000000).
        SetMinimumAmountOut(950000)

    instruction, err := swapInst.Build()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Swap instruction created with %d accounts\n", len(instruction.Accounts()))
    
    // Create a buy instruction
    buyInst := NewBuyInstruction().
        SetUserAuthority(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH")).
        SetTokenMint(solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")).
        SetAmount(1000000).
        SetMaxSolCost(500000)

    buyInstruction, err := buyInst.Build()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Buy instruction created with %d accounts\n", len(buyInstruction.Accounts()))
}
```

### Parsing Transactions

```go
package main

import (
    "fmt"
    "github.com/gagliardetto/solana-go"
)

func main() {
    // Parse a transaction from base64 encoded data
    txData := "your_base64_encoded_transaction_data"
    slot := uint64(123456789)

    transaction, err := ParseTransaction(txData, slot)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Transaction signature: %s\n", transaction.Signature)
    fmt.Printf("Number of trades: %d\n", len(transaction.Trade))
    fmt.Printf("Number of creates: %d\n", len(transaction.Create))
    fmt.Printf("Number of migrations: %d\n", len(transaction.Migrate))
}
```

## Available Instruction Builders

### SwapInstruction
```go
swapInst := NewSwapInstruction().
    SetUserSourceToken(sourceTokenAccount).
    SetUserDestToken(destTokenAccount).
    SetUserOwner(userWallet).
    SetAmmID(ammPool).
    SetAmountIn(1000000).
    SetMinimumAmountOut(950000)
```

### BuyInstruction
```go
buyInst := NewBuyInstruction().
    SetUserAuthority(userWallet).
    SetTokenMint(tokenMint).
    SetAmount(1000000).
    SetMaxSolCost(500000)
```

### SellInstruction
```go
sellInst := NewSellInstruction().
    SetUserAuthority(userWallet).
    SetTokenMint(tokenMint).
    SetAmount(1000000).
    SetMinSolReceived(450000)
```

### CreateTokenInstruction
```go
createInst := NewCreateTokenInstruction().
    SetPayer(payerWallet).
    SetMint(newTokenMint).
    SetDecimals(9).
    SetName("My Token").
    SetSymbol("MTK").
    SetInitialSupply(1000000000)
```

### MigrateInstruction
```go
migrateInst := NewMigrateInstruction().
    SetUserAuthority(userWallet).
    SetFromPool(oldPool).
    SetToPool(newPool).
    SetAmount(1000000)
```

## Testing

### Run All Tests
```bash
go test -v
```

### Run Specific Tests
```bash
go test -v -run TestSwapInstructionBuilder
go test -v -run TestBuyInstructionBuilder
go test -v -run TestSellInstructionBuilder
```

### Transaction Submission Tests
To test actual transaction submission, set environment variables:
```bash
export SOLANA_WALLET_PATH="/path/to/your/wallet.json"
export SOLANA_RPC_ENDPOINT="https://api.mainnet-beta.solana.com"
go test -v -run TestTransactionSubmission
```

## Supported Raydium Program IDs

- **Raydium V4**: `675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8`
- **Raydium V5**: `5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h`
- **Raydium Staking**: `EhhTKczWMGQt46ynNeRX1WfeagwwJd7ufHvCDjRxjo5Q`
- **Raydium Liquidity**: `27haf8L6oxUeXrHrgEgsexjSY5hbVUWEmvv9Nyxg8vQv`
- **Raydium Launchpad V1**: `6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P`
- **Raydium CP-Swap**: `CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C`

## Instruction Discriminators

| Operation | Discriminator | Description |
|-----------|---------------|-------------|
| Swap | 1 | Token swap through AMM |
| Buy | 6 | Token purchase in Launchpad |
| Sell | 7 | Token sale in Launchpad |
| Create Pool | 9 | Token/pool creation |
## Architecture

### Parser Module (`parser.go`)
- Handles both Geyser and standard RPC transaction formats
- Detects and parses various Raydium program IDs
- Implements generic parsing for unknown instruction discriminators
- Provides debug logging for instruction analysis

### Instruction Builders (`instructions.go`)
- Implements builder pattern for all major Raydium operations
- Provides fluent API with method chaining
- Handles serialization to valid Solana instructions
- Maintains separation between parsing and building functionality

### Type Definitions (`types.go`)
- Defines core transaction and instruction structures
- Provides type safety for all operations
- Supports extensibility for new instruction types

### Utilities (`utils.go`)
- Common helper functions for transaction analysis
- Validation and formatting utilities
- Error handling and logging support

## Error Handling

The library provides comprehensive error handling:
- Invalid transaction format detection
- Missing required fields validation
- Network connectivity issues
- RPC endpoint failures
- Invalid instruction data

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For issues and questions:
- Check the test files for usage examples
- Review the parser debug logs for transaction analysis
- Consult the Solana and Raydium documentation for protocol details
