package main

import (
	"encoding/hex"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// TokenInfo represents token metadata
type TokenInfo struct {
	Mint     solana.PublicKey
	Symbol   string
	Name     string
	Decimals uint8
}

// PoolInfo represents pool information
type PoolInfo struct {
	Address     solana.PublicKey
	TokenA      solana.PublicKey
	TokenB      solana.PublicKey
	TokenAVault solana.PublicKey
	TokenBVault solana.PublicKey
	LpMint      solana.PublicKey
	Fee         uint64
}

// Known token mints and their information
var knownTokens = map[string]TokenInfo{
	"So11111111111111111111111111111111111111112": {
		Mint:     solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112"),
		Symbol:   "SOL",
		Name:     "Solana",
		Decimals: 9,
	},
	"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": {
		Mint:     solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
		Symbol:   "USDC",
		Name:     "USD Coin",
		Decimals: 6,
	},
	"Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB": {
		Mint:     solana.MustPublicKeyFromBase58("Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"),
		Symbol:   "USDT",
		Name:     "Tether USD",
		Decimals: 6,
	},
}

// GetTokenInfo retrieves token information by mint address
func GetTokenInfo(mint solana.PublicKey) TokenInfo {
	if info, exists := knownTokens[mint.String()]; exists {
		return info
	}

	// Return default info for unknown tokens
	return TokenInfo{
		Mint:     mint,
		Symbol:   "UNKNOWN",
		Name:     "Unknown Token",
		Decimals: 9,
	}
}

// FormatTokenAmount formats a token amount according to its decimals
func FormatTokenAmount(amount uint64, decimals uint8) string {
	if decimals == 0 {
		return fmt.Sprintf("%d", amount)
	}

	divisor := uint64(1)
	for i := uint8(0); i < decimals; i++ {
		divisor *= 10
	}

	integerPart := amount / divisor
	fractionalPart := amount % divisor

	if fractionalPart == 0 {
		return fmt.Sprintf("%d", integerPart)
	}

	// Format with appropriate decimal places
	formatStr := fmt.Sprintf("%%d.%%0%dd", decimals)
	return fmt.Sprintf(formatStr, integerPart, fractionalPart)
}

// IsRaydiumProgram checks if a program ID is a known Raydium program
func IsRaydiumProgram(programID solana.PublicKey) bool {
	raydiumPrograms := []solana.PublicKey{
		RaydiumV4ProgramID,
		RaydiumV5ProgramID,
		RaydiumStakingProgramID,
		RaydiumLiquidityProgramID,
	}

	for _, program := range raydiumPrograms {
		if programID.Equals(program) {
			return true
		}
	}

	return false
}

// ExtractInstructionData extracts structured data from instruction bytes
func ExtractInstructionData(data []byte) map[string]interface{} {
	result := make(map[string]interface{})

	if len(data) == 0 {
		return result
	}

	result["discriminator"] = data[0]
	result["data_hex"] = hex.EncodeToString(data)
	result["data_length"] = len(data)

	// Try to extract common fields based on instruction format
	if len(data) >= 8 {
		// Assume next 8 bytes might be an amount (little-endian)
		amount := uint64(data[1]) | uint64(data[2])<<8 | uint64(data[3])<<16 | uint64(data[4])<<24 |
			uint64(data[5])<<32 | uint64(data[6])<<40 | uint64(data[7])<<48 | uint64(data[8])<<56
		result["potential_amount"] = amount
	}

	return result
}

// AnalyzeTransaction provides detailed analysis of a transaction
func AnalyzeTransaction(tx *Transaction) {
	fmt.Println("=== Transaction Analysis ===")

	// Analyze transaction type
	if len(tx.Create) > 0 {
		fmt.Println("🏗️  Pool/Token Creation Transaction")
	}
	if len(tx.Trade) > 0 {
		fmt.Println("💱 Trading Transaction")
	}
	if len(tx.Migrate) > 0 {
		fmt.Println("🔄 Migration Transaction")
	}

	// Analyze trading activity
	totalBuys := len(tx.SwapBuys)
	totalSells := len(tx.SwapSells)

	if totalBuys > 0 || totalSells > 0 {
		fmt.Printf("📊 Trading Activity: %d buys, %d sells\n", totalBuys, totalSells)
	}

	// Analyze tokens involved
	tokensInvolved := make(map[string]bool)
	for _, trade := range tx.Trade {
		tokensInvolved[trade.TokenIn.String()] = true
		tokensInvolved[trade.TokenOut.String()] = true
	}

	if len(tokensInvolved) > 0 {
		fmt.Printf("🪙 Tokens involved: %d unique tokens\n", len(tokensInvolved))
		for tokenAddr := range tokensInvolved {
			mint := solana.MustPublicKeyFromBase58(tokenAddr)
			tokenInfo := GetTokenInfo(mint)
			fmt.Printf("   - %s (%s)\n", tokenInfo.Symbol, tokenInfo.Name)
		}
	}

	fmt.Println()
}

// ValidateTransaction performs basic validation on a parsed transaction
func ValidateTransaction(tx *Transaction) []string {
	var issues []string

	// Check for zero signature
	if tx.Signature.IsZero() {
		issues = append(issues, "Transaction has zero signature")
	}

	// Check for reasonable slot number
	if tx.Slot == 0 {
		issues = append(issues, "Transaction has zero slot number")
	}

	// Validate trade consistency
	if len(tx.TradeBuys) != len(tx.SwapBuys) {
		issues = append(issues, "Mismatch between trade buys count and swap buys count")
	}

	if len(tx.TradeSells) != len(tx.SwapSells) {
		issues = append(issues, "Mismatch between trade sells count and swap sells count")
	}

	// Check for empty public keys in trades
	for i, trade := range tx.Trade {
		if trade.TokenIn.IsZero() || trade.TokenOut.IsZero() {
			issues = append(issues, fmt.Sprintf("Trade %d has zero token addresses", i))
		}
		if trade.Pool.IsZero() {
			issues = append(issues, fmt.Sprintf("Trade %d has zero pool address", i))
		}
	}

	return issues
}

// PrintValidationResults prints validation results
func PrintValidationResults(issues []string) {
	if len(issues) == 0 {
		fmt.Println("✅ Transaction validation passed")
		return
	}

	fmt.Printf("⚠️  Transaction validation found %d issues:\n", len(issues))
	for i, issue := range issues {
		fmt.Printf("   %d. %s\n", i+1, issue)
	}
}
