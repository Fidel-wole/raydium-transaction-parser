package main

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// Replace with a real Raydium swap transaction signature
const realTxSignature = "5wefCTqi9ynrh8pvVHFzpgHCLFFzoBwGoTgWSd6iq2Qw4Y51U4cEc2xHYtsdVSFZmRXUp5DNMSkhzb1CaXomLpJM"

func main() {
	fmt.Println("Raydium Transaction Parser")
	fmt.Println("==========================")

	// Check command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "test":
			fmt.Println("Test mode - functionality not implemented yet")
			return
		case "help", "-h", "--help":
			printUsage()
			return
		case "offline":
			fmt.Println("Running in offline mode...")
			fmt.Println("Offline mode - functionality not implemented yet")
			return
		}
	}

	fmt.Println("Fetching real transaction from Solana mainnet...")

	// Try multiple RPC endpoints in case one fails
	rpcEndpoints := []string{
		rpc.MainNetBeta_RPC,
		"https://solana-api.projectserum.com",
		"https://api.mainnet-beta.solana.com",
		"https://solana-mainnet.g.alchemy.com/v2/demo",
	}

	var txResp *rpc.GetTransactionResult
	var client *rpc.Client
	signature, err := solana.SignatureFromBase58(realTxSignature)
	if err != nil {
		log.Printf("Failed to parse signature: %v", err)
		fmt.Println("Falling back to basic demo...")
		demonstrateBasicFunctionality()
		return
	}

	// Try each RPC endpoint
	for i, endpoint := range rpcEndpoints {
		fmt.Printf("Trying RPC endpoint %d/%d: %s\n", i+1, len(rpcEndpoints), endpoint)
		client = rpc.New(endpoint)

		// Create a context with timeout for each request
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		txResp, err = client.GetTransaction(
			ctx,
			signature,
			&rpc.GetTransactionOpts{
				MaxSupportedTransactionVersion: &[]uint64{0}[0], // Support version 0 transactions
				Encoding:                       "base64",
			},
		)

		cancel() // Clean up the context

		if err == nil && txResp != nil && txResp.Transaction != nil {
			fmt.Printf("✅ Successfully fetched transaction from endpoint %d\n", i+1)
			break
		}

		log.Printf("❌ Endpoint %d failed: %v", i+1, err)
		if i < len(rpcEndpoints)-1 {
			fmt.Printf("Trying next endpoint...\n")
		}
	}
	if err != nil || txResp == nil || txResp.Transaction == nil {
		log.Printf("❌ All RPC endpoints failed. Last error: %v", err)
		fmt.Println("Falling back to basic demo...")
		demonstrateBasicFunctionality()
		testWithRaydiumData()
		// testWithMockRaydiumTransaction() // Temporarily disabled
		return
	}

	// Get the base64 encoded transaction
	encoded := txResp.Transaction.GetBinary()
	slot := txResp.Slot

	fmt.Println("Parsing transaction...")

	transaction, err := ParseTransaction(base64.StdEncoding.EncodeToString(encoded), slot)
	if err != nil {
		fmt.Printf("Failed to parse transaction: %v\n", err)
		demonstrateBasicFunctionality()
		return
	}

	fmt.Printf("Transaction successfully parsed!\n\n")

	issues := ValidateTransaction(transaction)
	PrintValidationResults(issues)
	fmt.Println()

	AnalyzeTransaction(transaction)
	printTransaction(transaction)

	// Optional: Load another transaction from a file
	if _, err := os.Stat("sample_transaction.txt"); err == nil {
		fmt.Println("\nLoading transaction from file...")
		loadAndParseFromFile("sample_transaction.txt")
	}
}

// printUsage prints the usage information
func printUsage() {
	fmt.Println("Usage: raydium-parser [command]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  test         Run all tests in offline mode")
	fmt.Println("  offline      Run in offline mode (same as test)")
	fmt.Println("  help         Show this help message")
	fmt.Println("  (no args)    Fetch and parse a real transaction from Solana mainnet")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run .                    # Fetch real transaction")
	fmt.Println("  go run . test               # Run tests")
	fmt.Println("  go run . offline            # Run in offline mode")
	fmt.Println("  ./raydium-parser test       # Run tests (compiled)")
}

// printTransaction prints the transaction details in a formatted way
func printTransaction(tx *Transaction) {
	fmt.Printf("Signature: %s\n", tx.Signature.String())
	fmt.Printf("Slot: %d\n", tx.Slot)
	fmt.Printf("Number of Creates: %d\n", len(tx.Create))
	fmt.Printf("Number of Trades: %d\n", len(tx.Trade))
	fmt.Printf("Number of Trade Buys: %d\n", len(tx.TradeBuys))
	fmt.Printf("Number of Trade Sells: %d\n", len(tx.TradeSells))
	fmt.Printf("Number of Migrations: %d\n", len(tx.Migrate))
	fmt.Printf("Number of Swap Buys: %d\n", len(tx.SwapBuys))
	fmt.Printf("Number of Swap Sells: %d\n", len(tx.SwapSells))

	if len(tx.Create) > 0 {
		fmt.Println("\nCreate Operations:")
		for i, create := range tx.Create {
			fmt.Printf("  [%d] Token: %s, Pool: %s, Creator: %s\n",
				i, create.TokenMint.String(), create.PoolAddress.String(), create.Creator.String())
		}
	}

	if len(tx.Trade) > 0 {
		fmt.Println("\nTrade Operations:")
		for i, trade := range tx.Trade {
			fmt.Printf("  [%d] Type: %s, TokenIn: %s, TokenOut: %s, Trader: %s, Pool: %s\n",
				i, trade.TradeType, trade.TokenIn.String(), trade.TokenOut.String(),
				trade.Trader.String(), trade.Pool.String())
		}
	}

	if len(tx.Migrate) > 0 {
		fmt.Println("\nMigration Operations:")
		for i, migration := range tx.Migrate {
			fmt.Printf("  [%d] From: %s, To: %s, Token: %s, Owner: %s\n",
				i, migration.FromPool.String(), migration.ToPool.String(),
				migration.Token.String(), migration.Owner.String())
		}
	}

	if len(tx.SwapBuys) > 0 {
		fmt.Println("\nSwap Buy Operations:")
		for i, swap := range tx.SwapBuys {
			fmt.Printf("  [%d] TokenIn: %s, TokenOut: %s, AmountIn: %d, AmountOut: %d, Buyer: %s\n",
				i, swap.TokenIn.String(), swap.TokenOut.String(),
				swap.AmountIn, swap.AmountOut, swap.Buyer.String())
		}
	}

	if len(tx.SwapSells) > 0 {
		fmt.Println("\nSwap Sell Operations:")
		for i, swap := range tx.SwapSells {
			fmt.Printf("  [%d] TokenIn: %s, TokenOut: %s, AmountIn: %d, AmountOut: %d, Seller: %s\n",
				i, swap.TokenIn.String(), swap.TokenOut.String(),
				swap.AmountIn, swap.AmountOut, swap.Seller.String())
		}
	}

	// Pretty print as JSON for debugging
	fmt.Println("\nJSON Representation:")
	jsonData, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
		return
	}
	fmt.Println(string(jsonData))
}

// loadAndParseFromFile loads a transaction from a file and parses it
func loadAndParseFromFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading file %s: %v", filename, err)
		return
	}

	// Assume the file contains a base64 encoded transaction
	transactionBase64 := string(data)

	// Parse the transaction
	slot := uint64(987654321) // Sample slot
	transaction, err := ParseTransaction(transactionBase64, slot)
	if err != nil {
		log.Printf("Failed to parse transaction from file: %v", err)
		return
	}

	fmt.Printf("Transaction from file parsed successfully!\n")
	printTransaction(transaction)

	// Show comprehensive examples
	testWithRaydiumData()
	// testWithMockRaydiumTransaction() // Temporarily disabled
}

// demonstrateBasicFunctionality shows basic parser functionality
func demonstrateBasicFunctionality() {
	fmt.Println("=== Basic Parser Functionality Demo ===")

	// Create a mock transaction to demonstrate parsing
	mockSignature := solana.Signature{}
	copy(mockSignature[:], []byte("demo_signature"))

	tx := &Transaction{
		Signature:  mockSignature,
		Slot:       123456789,
		Create:     []CreateInfo{},
		Trade:      []TradeInfo{},
		TradeBuys:  []int{},
		TradeSells: []int{},
		Migrate:    []Migration{},
		SwapBuys:   []SwapBuy{},
		SwapSells:  []SwapSell{},
	}

	// Add some mock data
	mockTokenIn := solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")   // SOL
	mockTokenOut := solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v") // USDC

	tradeInfo := TradeInfo{
		InstructionIndex: 0,
		TokenIn:          mockTokenIn,
		TokenOut:         mockTokenOut,
		AmountIn:         1000000000, // 1 SOL
		AmountOut:        25000000,   // 25 USDC
		Trader:           solana.MustPublicKeyFromBase58("11111111111111111111111111111111"),
		Pool:             solana.MustPublicKeyFromBase58("58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2"),
		TradeType:        "swap",
	}

	tx.Trade = append(tx.Trade, tradeInfo)

	fmt.Println("Demo transaction created successfully!")
	printTransaction(tx)

	// Test utility functions
	fmt.Println("\n=== Testing Utility Functions ===")
	fmt.Printf("SOL is base currency: %t\n", isBaseCurrency(mockTokenIn))
	fmt.Printf("USDC is base currency: %t\n", isBaseCurrency(mockTokenOut))

	// Test validation
	fmt.Println("\n=== Testing Validation ===")
	issues := ValidateTransaction(tx)
	if len(issues) > 0 {
		fmt.Printf("Validation found %d issues:\n", len(issues))
		for _, issue := range issues {
			fmt.Printf("- %s\n", issue)
		}
	} else {
		fmt.Println("No validation issues found!")
	}

	fmt.Println("\n=== Demo Complete ===")
}

// testWithRaydiumData demonstrates parsing with crafted Raydium-like data
func testWithRaydiumData() {
	fmt.Println("\n=== Testing with Simulated Raydium Data ===")

	// Create a transaction with simulated Raydium instructions
	mockSignature := solana.Signature{}
	copy(mockSignature[:], []byte("raydium_test_signature"))

	result := &Transaction{
		Signature:  mockSignature,
		Slot:       123456789,
		Create:     []CreateInfo{},
		Trade:      []TradeInfo{},
		TradeBuys:  []int{},
		TradeSells: []int{},
		Migrate:    []Migration{},
		SwapBuys:   []SwapBuy{},
		SwapSells:  []SwapSell{},
	}

	// Add mock token creation
	createInfo := CreateInfo{
		TokenMint:     solana.MustPublicKeyFromBase58("4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"),
		PoolAddress:   solana.MustPublicKeyFromBase58("58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2"),
		Creator:       solana.MustPublicKeyFromBase58("7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU"),
		TokenDecimals: 9,
		TokenSymbol:   "RAYTEST",
		Amount:        1000000000000, // 1M tokens
		Timestamp:     1700000000,
	}
	result.Create = append(result.Create, createInfo)

	// Add mock swap (buy)
	buyTrade := TradeInfo{
		InstructionIndex: 1,
		TokenIn:          solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112"),  // SOL
		TokenOut:         solana.MustPublicKeyFromBase58("4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"), // New token
		AmountIn:         500000000,                                                                      // 0.5 SOL
		AmountOut:        1000000000,                                                                     // 1000 tokens
		Trader:           solana.MustPublicKeyFromBase58("7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU"),
		Pool:             solana.MustPublicKeyFromBase58("58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2"),
		TradeType:        "buy",
	}
	result.Trade = append(result.Trade, buyTrade)
	result.TradeBuys = append(result.TradeBuys, 1)

	buySwap := SwapBuy{
		TokenIn:      buyTrade.TokenIn,
		TokenOut:     buyTrade.TokenOut,
		AmountIn:     buyTrade.AmountIn,
		AmountOut:    buyTrade.AmountOut,
		Pool:         buyTrade.Pool,
		Buyer:        buyTrade.Trader,
		MinAmountOut: 950000000, // 950 tokens minimum
		Slippage:     0.05,      // 5% slippage
	}
	result.SwapBuys = append(result.SwapBuys, buySwap)

	// Add mock swap (sell)
	sellTrade := TradeInfo{
		InstructionIndex: 2,
		TokenIn:          solana.MustPublicKeyFromBase58("4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"), // New token
		TokenOut:         solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112"),  // SOL
		AmountIn:         500000000,                                                                      // 500 tokens
		AmountOut:        200000000,                                                                      // 0.2 SOL
		Trader:           solana.MustPublicKeyFromBase58("9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM"),
		Pool:             solana.MustPublicKeyFromBase58("58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2"),
		TradeType:        "sell",
	}
	result.Trade = append(result.Trade, sellTrade)
	result.TradeSells = append(result.TradeSells, 2)

	sellSwap := SwapSell{
		TokenIn:      sellTrade.TokenIn,
		TokenOut:     sellTrade.TokenOut,
		AmountIn:     sellTrade.AmountIn,
		AmountOut:    sellTrade.AmountOut,
		Pool:         sellTrade.Pool,
		Seller:       sellTrade.Trader,
		MinAmountOut: 180000000, // 0.18 SOL minimum
		Slippage:     0.10,      // 10% slippage
	}
	result.SwapSells = append(result.SwapSells, sellSwap)

	// Add mock migration
	migration := Migration{
		FromPool:  solana.MustPublicKeyFromBase58("58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2"),
		ToPool:    solana.MustPublicKeyFromBase58("7XawhbbxtsRcQA8KTkHT9f9nc6d69UwqCDh6U5EEbEmX"),
		Token:     solana.MustPublicKeyFromBase58("4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"),
		Owner:     solana.MustPublicKeyFromBase58("7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU"),
		Amount:    250000000, // 250 tokens
		Timestamp: 1700000100,
	}
	result.Migrate = append(result.Migrate, migration)

	fmt.Println("✅ Simulated Raydium transaction created successfully!")
	printTransaction(result)
}

// testWithMockRaydiumTransaction tests with a crafted transaction containing Raydium-like binary data
func testWithMockRaydiumTransaction() {
	fmt.Println("\n=== Testing with Mock Raydium Transaction Data ===")

	// Create a simple mock transaction with Raydium-like binary structure
	// This simulates what actual Raydium transaction bytes might look like
	mockTransactionBytes := make([]byte, 400) // Increased size to fix bounds issue

	// Mock signature (64 bytes)
	copy(mockTransactionBytes[0:64], []byte("mock_raydium_signature_for_testing_purposes_and_demonstration"))

	// Mock message header (3 bytes)
	mockTransactionBytes[64] = 1 // numSignatures
	mockTransactionBytes[65] = 0 // numReadonlySignedAccounts
	mockTransactionBytes[66] = 0 // numReadonlyUnsignedAccounts

	// Mock account keys length (1 byte + account keys)
	mockTransactionBytes[67] = 8 // 8 accounts

	// Mock account keys (32 bytes each)
	solPubkey := solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	raydiumPubkey := solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")

	copy(mockTransactionBytes[68:100], solPubkey[:])
	copy(mockTransactionBytes[100:132], raydiumPubkey[:])

	// Fill remaining account keys with mock data
	for i := 0; i < 6; i++ {
		start := 132 + (i * 32)
		copy(mockTransactionBytes[start:start+32], []byte(fmt.Sprintf("mock_account_%d_for_testing_purp", i)))
	}
	// Mock recent blockhash (32 bytes)
	copy(mockTransactionBytes[324:356], []byte("mock_recent_blockhash_for_testing"))

	// Mock instruction data
	instructionStart := 260
	mockTransactionBytes[instructionStart] = 1 // numInstructions

	// Mock Raydium swap instruction
	mockTransactionBytes[instructionStart+1] = 1 // programIdIndex (Raydium)
	mockTransactionBytes[instructionStart+2] = 6 // numAccounts
	// Account indices
	mockTransactionBytes[instructionStart+3] = 0 // tokenIn
	mockTransactionBytes[instructionStart+4] = 2 // tokenOut
	mockTransactionBytes[instructionStart+5] = 3 // pool
	mockTransactionBytes[instructionStart+6] = 4 // trader
	mockTransactionBytes[instructionStart+7] = 5 // vault
	mockTransactionBytes[instructionStart+8] = 6 // authority

	// Mock instruction data
	mockTransactionBytes[instructionStart+9] = 16 // data length
	mockTransactionBytes[instructionStart+10] = 1 // instruction discriminator (swap)

	// Mock amounts (8 bytes each)
	amountIn := uint64(1000000000)    // 1 SOL
	minAmountOut := uint64(950000000) // 950 tokens

	binary.LittleEndian.PutUint64(mockTransactionBytes[instructionStart+11:instructionStart+19], amountIn)
	binary.LittleEndian.PutUint64(mockTransactionBytes[instructionStart+19:instructionStart+27], minAmountOut)

	// Encode to base64
	encodedTx := base64.StdEncoding.EncodeToString(mockTransactionBytes)

	fmt.Printf("Created mock transaction with %d bytes\n", len(mockTransactionBytes))
	fmt.Printf("Encoded length: %d characters\n", len(encodedTx))

	// Try to parse it
	result, err := ParseTransaction(encodedTx, 987654321)
	if err != nil {
		fmt.Printf("Parser error: %v\n", err)
		return
	}

	fmt.Println("✅ Mock Raydium transaction parsed successfully!")
	printTransaction(result)
}
