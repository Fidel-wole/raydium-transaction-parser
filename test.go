package main

import (
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
)

// TestTransactionParsing tests the basic transaction parsing functionality
func TestTransactionParsing() {
	// Test with empty transaction
	_, err := ParseTransaction("", 12345)
	if err == nil {
		log.Println("Expected error for empty transaction")
		return
	}

	// Test with sample transaction
	sampleTx := "VGVzdCB0cmFuc2FjdGlvbiBkYXRhIGZvciBSYXlkaXVtIHBhcnNlciB0ZXN0aW5nIHB1cnBvc2Vz"
	result, err := ParseTransaction(sampleTx, 12345)
	if err != nil {
		log.Printf("Error parsing sample transaction: %v", err)
		return
	}

	if result.Slot != 12345 {
		log.Printf("Expected slot 12345, got %d", result.Slot)
		return
	}

	fmt.Println("✓ Basic transaction parsing tests passed")
}

// TestTokenUtilities tests token-related utility functions
func TestTokenUtilities() {
	// Test known token lookup
	solMint := solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	
	if !isBaseCurrency(solMint) {
		log.Printf("SOL should be recognized as base currency")
		return
	}

	// Test unknown token
	randomMint := solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
	if isBaseCurrency(randomMint) {
		log.Printf("Random mint should not be base currency")
		return
	}

	fmt.Println("✓ Token utility tests passed")
}

// TestValidation tests transaction validation
func TestValidation() {
	// Create a mock transaction
	mockSignature := solana.Signature{}
	copy(mockSignature[:], []byte("test_signature"))

	tx := &Transaction{
		Signature:  mockSignature,
		Slot:       12345,
		Create:     []CreateInfo{},
		Trade:      []TradeInfo{},
		TradeBuys:  []int{},
		TradeSells: []int{},
		Migrate:    []Migration{},
		SwapBuys:   []SwapBuy{},
		SwapSells:  []SwapSell{},
	}

	issues := ValidateTransaction(tx)
	fmt.Printf("✓ Validation completed with %d issues\n", len(issues))
}

// TestInstructionParsing tests instruction parsing utilities
func TestInstructionParsing() {
	// Test program ID recognition
	if !isRaydiumProgram(RaydiumV4ProgramID) {
		log.Printf("Raydium V4 program ID should be recognized")
		return
	}

	if !isRaydiumProgram(RaydiumV5ProgramID) {
		log.Printf("Raydium V5 program ID should be recognized")
		return
	}

	// Test base currency detection
	solMint := solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	if !isBaseCurrency(solMint) {
		log.Printf("SOL should be recognized as base currency")
		return
	}

	randomMint := solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
	if isBaseCurrency(randomMint) {
		log.Printf("Random mint should not be base currency")
		return
	}

	fmt.Println("✓ Instruction parsing tests passed")
}

// isRaydiumProgram checks if a program ID is a known Raydium program
func isRaydiumProgram(programID solana.PublicKey) bool {
	return programID.Equals(RaydiumV4ProgramID) ||
		programID.Equals(RaydiumV5ProgramID) ||
		programID.Equals(RaydiumStakingProgramID) ||
		programID.Equals(RaydiumLiquidityProgramID)
}

// RunAllTests runs all test functions
func RunAllTests() {
	fmt.Println("Running Raydium Parser Tests")
	fmt.Println("============================")

	TestTransactionParsing()
	TestTokenUtilities()
	TestValidation()
	TestInstructionParsing()

	fmt.Println("\n✓ All tests completed!")
}
