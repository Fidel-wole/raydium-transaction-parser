# Raydium Transaction Parser

A Go project for parsing Solana transactions with a focus on Raydium DEX operations.

## Features

- Parse Solana transactions using the `solana-go` library
- Detect and extract Raydium-specific operations:
  - Token/Pool creation
  - Swaps (buy/sell)
  - Liquidity operations
  - Migrations
- Support for both top-level and inner instructions
- Structured data extraction for trading analysis

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

### Running the Parser

```bash
go run .
```

This will parse a sample transaction and display the results.

### Custom Transaction Parsing

You can parse your own transactions by:

1. Creating a file named `sample_transaction.txt` with a base64-encoded transaction
2. Or modifying the `sampleTransactionBase64` constant in `main.go`

### Example Output

```
Raydium Transaction Parser
==========================
Parsing sample transaction...
Transaction successfully parsed!

Signature: [signature]
Slot: 123456789
Number of Creates: 0
Number of Trades: 0
Number of Trade Buys: 0
Number of Trade Sells: 0
Number of Migrations: 0
Number of Swap Buys: 0
Number of Swap Sells: 0
```

## Project Structure

- `main.go` - Entry point and sample usage
- `parser.go` - Core parsing logic and instruction handlers
- `types.go` - Data structures for parsed transaction data
- `go.mod` - Go module definition
- `README.md` - This file

## Data Structures

### Transaction
The main structure containing all parsed transaction data:
- `Signature` - Transaction signature
- `Slot` - Block slot number
- `Create` - Token/pool creation operations
- `Trade` - General trade information
- `TradeBuys/TradeSells` - Buy/sell operation indices
- `Migrate` - Migration operations
- `SwapBuys/SwapSells` - Detailed swap information

### Supporting Types
- `CreateInfo` - Token/pool creation details
- `TradeInfo` - Trade operation details
- `Migration` - Migration operation details
- `SwapBuy/SwapSell` - Detailed swap operation data

## Known Raydium Program IDs

The parser recognizes the following Raydium program IDs:
- V4: `675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8`
- V5: `5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h`
- Staking: `EhhTKczWMGQt46ynNeRX1WfeagwwJd7ufHvCDjRxjo5Q`
- Liquidity: `27haf8L6oxUeXrHrgEgsexjSY5hbVUWEmvv9Nyxg8vQv`

## Future Enhancements

1. **Inner Instruction Parsing**: Currently only parses top-level instructions
2. **Real Transaction Data**: Integration with actual Raydium transactions
3. **Token Metadata**: Fetch token symbols and metadata
4. **Price Calculation**: Calculate USD values for trades
5. **Geyser Integration**: Connect to Solana Geyser plugin for real-time data
6. **Database Storage**: Store parsed data in a database
7. **API Endpoints**: REST API for querying parsed transactions

## Development Notes

- The current implementation uses placeholder instruction discriminators
- Actual parsing logic would need to be adapted based on the real Raydium IDL
- Some fields are currently hardcoded and would need dynamic extraction
- Error handling could be enhanced for production use

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is provided as-is for educational and development purposes.
# raydium-transaction-parser
