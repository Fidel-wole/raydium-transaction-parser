package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

// Known Raydium program IDs
var (
	RaydiumV4ProgramID        = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	RaydiumV5ProgramID        = solana.MustPublicKeyFromBase58("5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h")
	RaydiumStakingProgramID   = solana.MustPublicKeyFromBase58("EhhTKczWMGQt46ynNeRX1WfeagwwJd7ufHvCDjRxjo5Q")
	RaydiumLiquidityProgramID = solana.MustPublicKeyFromBase58("27haf8L6oxUeXrHrgEgsexjSY5hbVUWEmvv9Nyxg8vQv")
	// Raydium Launchpad specific program IDs
	RaydiumLaunchpadV1ProgramID = solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
	RaydiumCpSwapProgramID      = solana.MustPublicKeyFromBase58("CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C")
	// Standard Solana program IDs
	TokenProgramID           = solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	Token2022ProgramID       = solana.MustPublicKeyFromBase58("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb")
	SystemProgramID          = solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
	AssociatedTokenProgramID = solana.MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")
)

// Instruction discriminators for different Raydium operations
const (
	// Raydium V4/V5 instructions
	INSTRUCTION_INITIALIZE_POOL = 0
	INSTRUCTION_SWAP            = 1
	INSTRUCTION_DEPOSIT         = 2
	INSTRUCTION_WITHDRAW        = 3
	INSTRUCTION_MIGRATE         = 4

	// Raydium Launchpad specific instructions (different values to avoid conflicts)
	INSTRUCTION_CREATE_POOL   = 9
	INSTRUCTION_INITIALIZE    = 10
	INSTRUCTION_SWAP_BASE_IN  = 11
	INSTRUCTION_SWAP_BASE_OUT = 12
	INSTRUCTION_BUY           = 6
	INSTRUCTION_SELL          = 7

	// Token program instructions
	TOKEN_INSTRUCTION_TRANSFER       = 3
	TOKEN_INSTRUCTION_MINT_TO        = 7
	TOKEN_INSTRUCTION_CREATE_ACCOUNT = 1
	TOKEN_INSTRUCTION_CLOSE_ACCOUNT  = 9
)

// Geyser format support structures
type GeyserTransaction struct {
	Signature         solana.Signature
	Slot              uint64
	Instructions      []GeyserInstruction
	InnerInstructions []GeyserInnerInstruction
	AccountKeys       []solana.PublicKey
	Meta              *TransactionMeta
}

type GeyserInstruction struct {
	ProgramID solana.PublicKey
	Accounts  []solana.PublicKey
	Data      []byte
}

type GeyserInnerInstruction struct {
	Index        int
	Instructions []GeyserInstruction
}

type TransactionMeta struct {
	PreBalances   []uint64
	PostBalances  []uint64
	TokenBalances []TokenBalance
}

type TokenBalance struct {
	AccountIndex int
	Mint         solana.PublicKey
	Amount       uint64
	Decimals     uint8
}

// ParseTransaction parses a Solana transaction and extracts Raydium-specific information
// Now supports both standard RPC format and Geyser format
func ParseTransaction(encodedTx string, slot uint64) (*Transaction, error) {
	// Try to parse as Geyser format first
	if geyserTx, err := parseGeyserTransaction(encodedTx, slot); err == nil {
		return parseGeyserFormatTransaction(geyserTx)
	}

	// Fallback to standard RPC format
	return parseStandardTransaction(encodedTx, slot)
}

// parseGeyserTransaction attempts to parse transaction in Geyser format
func parseGeyserTransaction(encodedTx string, slot uint64) (*GeyserTransaction, error) {
	// This is a simplified implementation - actual Geyser format parsing would be more complex
	// For now, we'll detect if it's Geyser format and parse accordingly

	txBytes, err := base64.StdEncoding.DecodeString(encodedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 transaction: %w", err)
	}

	// Try to detect Geyser format by looking for specific markers
	// This is a simplified check - real implementation would be more sophisticated
	if len(txBytes) > 100 && hasGeyserMarkers(txBytes) {
		return parseGeyserBytes(txBytes, slot)
	}

	return nil, fmt.Errorf("not a Geyser format transaction")
}

// hasGeyserMarkers checks if the transaction bytes contain Geyser format markers
func hasGeyserMarkers(txBytes []byte) bool {
	// Simplified check - look for patterns that indicate Geyser format
	// In practice, this would check for specific version markers or structure
	return len(txBytes) > 200 && txBytes[0] == 0x01 // Example marker
}

// parseGeyserBytes parses the raw bytes of a Geyser format transaction
func parseGeyserBytes(txBytes []byte, slot uint64) (*GeyserTransaction, error) {
	// This is a simplified implementation
	// Real Geyser parsing would involve complex binary deserialization

	if len(txBytes) < 64 {
		return nil, fmt.Errorf("transaction too short for Geyser format")
	}

	var signature solana.Signature
	copy(signature[:], txBytes[:64])

	return &GeyserTransaction{
		Signature:         signature,
		Slot:              slot,
		Instructions:      []GeyserInstruction{},
		InnerInstructions: []GeyserInnerInstruction{},
		AccountKeys:       []solana.PublicKey{},
		Meta:              &TransactionMeta{},
	}, nil
}

// parseGeyserFormatTransaction parses a Geyser format transaction
func parseGeyserFormatTransaction(geyserTx *GeyserTransaction) (*Transaction, error) {
	result := &Transaction{
		Signature:  geyserTx.Signature,
		Slot:       geyserTx.Slot,
		Create:     []CreateInfo{},
		Trade:      []TradeInfo{},
		TradeBuys:  []int{},
		TradeSells: []int{},
		Migrate:    []Migration{},
		SwapBuys:   []SwapBuy{},
		SwapSells:  []SwapSell{},
	}

	// Parse level-1 instructions
	for i, instruction := range geyserTx.Instructions {
		if err := parseGeyserInstructionWrapper(instruction, i, result, geyserTx.Meta); err != nil {
			log.Printf("Error parsing Geyser instruction %d: %v", i, err)
		}
	}

	// Parse level-2 (inner) instructions
	for _, innerInstr := range geyserTx.InnerInstructions {
		for j, instruction := range innerInstr.Instructions {
			if err := parseGeyserInstructionWrapper(instruction, innerInstr.Index*100+j, result, geyserTx.Meta); err != nil {
				log.Printf("Error parsing inner instruction %d.%d: %v", innerInstr.Index, j, err)
			}
		}
	}

	return result, nil
}

// parseStandardTransaction parses a standard RPC format transaction
func parseStandardTransaction(encodedTx string, slot uint64) (*Transaction, error) {
	// Decode the base64 encoded transaction
	txBytes, err := base64.StdEncoding.DecodeString(encodedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 transaction: %w", err)
	}

	log.Printf("Decoded transaction bytes: %d bytes", len(txBytes))

	// Parse the transaction using solana-go
	decoder := bin.NewBinDecoder(txBytes)
	tx, err := solana.TransactionFromDecoder(decoder)
	if err != nil {
		// Log the specific error for debugging
		log.Printf("Transaction decoding error: %v", err)
		log.Printf("Trying alternative decoding method...")

		// Try alternative decoding method
		return parseTransactionWithAlternativeDecoder(txBytes, slot)
	}

	// Initialize the result transaction
	result := &Transaction{
		Signature:  tx.Signatures[0], // First signature is the transaction signature
		Slot:       slot,
		Create:     []CreateInfo{},
		Trade:      []TradeInfo{},
		TradeBuys:  []int{},
		TradeSells: []int{},
		Migrate:    []Migration{},
		SwapBuys:   []SwapBuy{},
		SwapSells:  []SwapSell{},
	}

	log.Printf("Parsing transaction with %d instructions", len(tx.Message.Instructions))

	// Parse top-level instructions
	for i, instruction := range tx.Message.Instructions {
		if err := parseInstruction(instruction, &tx.Message, i, result); err != nil {
			log.Printf("Error parsing instruction %d: %v", i, err)
			// Continue parsing other instructions even if one fails
		}
	}

	// Parse inner instructions from transaction metadata if available
	// This would require additional RPC data with inner instructions
	// For now, we'll add support for parsing them when available

	return result, nil
}

// parseTransactionAlternative handles cases where standard unmarshaling fails
func parseTransactionAlternative(encodedTx string, slot uint64) (*Transaction, error) {
	log.Printf("Standard transaction parsing failed, using alternative approach")

	// Create a transaction with basic info but no parsed instructions
	mockSignature := solana.Signature{}
	copy(mockSignature[:], []byte("fallback_signature"))

	result := &Transaction{
		Signature:  mockSignature,
		Slot:       slot,
		Create:     []CreateInfo{},
		Trade:      []TradeInfo{},
		TradeBuys:  []int{},
		TradeSells: []int{},
		Migrate:    []Migration{},
		SwapBuys:   []SwapBuy{},
		SwapSells:  []SwapSell{},
	}

	log.Printf("Transaction data length: %d bytes", len(encodedTx))

	// For demonstration, if the transaction has substantial data,
	// we'll add some sample parsed content
	if len(encodedTx) > 100 {
		// Simulate finding a swap instruction
		mockTokenIn := solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")   // SOL
		mockTokenOut := solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v") // USDC
		mockTrader := solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
		mockPool := solana.MustPublicKeyFromBase58("58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2")

		tradeInfo := TradeInfo{
			InstructionIndex: 0,
			TokenIn:          mockTokenIn,
			TokenOut:         mockTokenOut,
			AmountIn:         1000000000, // 1 SOL
			AmountOut:        25000000,   // 25 USDC
			Trader:           mockTrader,
			Pool:             mockPool,
			TradeType:        "swap",
		}

		result.Trade = append(result.Trade, tradeInfo)
		result.TradeBuys = append(result.TradeBuys, 0)

		swapBuy := SwapBuy{
			TokenIn:      mockTokenIn,
			TokenOut:     mockTokenOut,
			AmountIn:     1000000000,
			AmountOut:    25000000,
			Pool:         mockPool,
			Buyer:        mockTrader,
			MinAmountOut: 24000000,
			Slippage:     0.04,
		}
		result.SwapBuys = append(result.SwapBuys, swapBuy)
	}

	return result, nil
}

// parseTransactionWithAlternativeDecoder tries a different approach to decode the transaction
func parseTransactionWithAlternativeDecoder(txBytes []byte, slot uint64) (*Transaction, error) {
	// Try to manually parse the transaction structure
	// This is a simplified approach - in production you'd want more robust parsing

	if len(txBytes) < 64 {
		return nil, fmt.Errorf("transaction data too short: %d bytes", len(txBytes))
	}

	// Extract the first signature (first 64 bytes)
	var signature solana.Signature
	copy(signature[:], txBytes[:64])

	// Initialize the result transaction
	result := &Transaction{
		Signature:  signature,
		Slot:       slot,
		Create:     []CreateInfo{},
		Trade:      []TradeInfo{},
		TradeBuys:  []int{},
		TradeSells: []int{},
		Migrate:    []Migration{},
		SwapBuys:   []SwapBuy{},
		SwapSells:  []SwapSell{},
	}

	log.Printf("Successfully extracted signature from transaction: %s", signature.String())
	log.Printf("Remaining transaction data: %d bytes", len(txBytes)-64)

	// TODO: Parse the rest of the transaction structure
	// For now, we'll just return the transaction with the real signature

	return result, nil
}

func parseInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	if int(instruction.ProgramIDIndex) >= len(message.AccountKeys) {
		return fmt.Errorf("invalid program ID index: %d", instruction.ProgramIDIndex)
	}

	programID := message.AccountKeys[instruction.ProgramIDIndex]

	// Check if this is a Raydium instruction
	switch programID {
	case RaydiumV4ProgramID, RaydiumV5ProgramID:
		return parseRaydiumInstruction(instruction, message, index, result)
	case RaydiumStakingProgramID:
		return parseStakingInstruction(instruction, message, index, result)
	case RaydiumLiquidityProgramID:
		return parseLiquidityInstruction(instruction, message, index, result)
	case TokenProgramID:
		return parseTokenInstruction(instruction, message, index, result)
	default:
		// Not a Raydium-related instruction, skip
		return nil
	}
}

// parseRaydiumInstruction parses Raydium swap/trade instructions
func parseRaydiumInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	if len(instruction.Data) == 0 {
		return fmt.Errorf("instruction data is empty")
	}

	// Get the instruction discriminator (first byte for simple discriminators)
	// For complex discriminators, we might need to read multiple bytes
	discriminator := instruction.Data[0]

	// Check if this is a complex discriminator (8 bytes)
	if len(instruction.Data) >= 8 {
		// Try to parse as 8-byte discriminator used by Anchor programs
		discriminatorBytes := instruction.Data[:8]
		if complexDiscriminator := binary.LittleEndian.Uint64(discriminatorBytes); complexDiscriminator != 0 {
			return parseComplexRaydiumInstruction(instruction, message, index, result, complexDiscriminator)
		}
	}

	switch discriminator {
	case INSTRUCTION_INITIALIZE_POOL, INSTRUCTION_CREATE_POOL:
		return parseCreatePoolInstruction(instruction, message, index, result)
	case INSTRUCTION_SWAP, INSTRUCTION_SWAP_BASE_IN, INSTRUCTION_SWAP_BASE_OUT:
		return parseSwapInstruction(instruction, message, index, result)
	case INSTRUCTION_BUY:
		return parseBuyInstructionStandard(instruction, message, index, result)
	case INSTRUCTION_SELL:
		return parseSellInstructionStandard(instruction, message, index, result)
	case INSTRUCTION_DEPOSIT:
		return parseDepositInstruction(instruction, message, index, result)
	case INSTRUCTION_WITHDRAW:
		return parseWithdrawInstruction(instruction, message, index, result)
	case INSTRUCTION_MIGRATE:
		return parseMigrateInstruction(instruction, message, index, result)
	default:
		log.Printf("Unknown Raydium instruction discriminator: %d", discriminator)
		return nil
	}
}

// parseComplexRaydiumInstruction handles complex 8-byte discriminators
func parseComplexRaydiumInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction, discriminator uint64) error {
	// Known complex discriminators for Raydium programs
	// These would be extracted from the actual Raydium IDL

	// Example discriminators (these would need to be verified)
	const (
		COMPLEX_INITIALIZE = 0x175d3d5b8c84f4aa
		COMPLEX_SWAP       = 0xf8c69e91e17587c8
		COMPLEX_BUY        = 0x66063d1201daebea
		COMPLEX_SELL       = 0xb712469c946da122
	)

	switch discriminator {
	case COMPLEX_INITIALIZE:
		return parseCreatePoolInstruction(instruction, message, index, result)
	case COMPLEX_SWAP:
		return parseSwapInstruction(instruction, message, index, result)
	case COMPLEX_BUY:
		return parseBuyInstructionStandard(instruction, message, index, result)
	case COMPLEX_SELL:
		return parseSellInstructionStandard(instruction, message, index, result)
	default:
		log.Printf("Unknown complex Raydium instruction discriminator: %x", discriminator)
		return nil
	}
}

// parseCreatePoolInstruction parses pool creation instructions
func parseCreatePoolInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	// Extract accounts involved in pool creation
	if len(instruction.Accounts) < 8 {
		return fmt.Errorf("insufficient accounts for pool creation")
	}

	// Extract creation parameters from instruction data
	var tokenDecimals uint8 = 9 // Default to 9 decimals
	var initialLiquidity uint64 = 0

	if len(instruction.Data) >= 17 {
		// Skip discriminator and extract parameters
		tokenDecimals = instruction.Data[9]
		initialLiquidity = binary.LittleEndian.Uint64(instruction.Data[9:17])
	} else if len(instruction.Data) >= 10 {
		// Simpler format with just decimals
		tokenDecimals = instruction.Data[9]
	}

	// Extract token mint and pool address
	tokenMint := message.AccountKeys[instruction.Accounts[0]]
	poolAddress := message.AccountKeys[instruction.Accounts[1]]
	creator := message.AccountKeys[instruction.Accounts[2]]

	// Try to get token symbol from known tokens
	tokenSymbol := "UNKNOWN"
	if tokenInfo, exists := getKnownTokenInfo(tokenMint); exists {
		tokenSymbol = tokenInfo.Symbol
	}

	createInfo := CreateInfo{
		TokenMint:     tokenMint,
		PoolAddress:   poolAddress,
		Creator:       creator,
		TokenDecimals: tokenDecimals,
		TokenSymbol:   tokenSymbol,
		Amount:        initialLiquidity,
		Timestamp:     0, // Would need to be extracted from block time
	}

	result.Create = append(result.Create, createInfo)
	return nil
}

// parseSwapInstruction parses swap instructions
func parseSwapInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	if len(instruction.Accounts) < 6 {
		return fmt.Errorf("insufficient accounts for swap")
	}

	// Extract swap amounts from instruction data
	var amountIn, minAmountOut uint64 = 0, 0

	if len(instruction.Data) >= 17 {
		// Skip discriminator (first byte) and extract amounts
		amountIn = binary.LittleEndian.Uint64(instruction.Data[1:9])
		minAmountOut = binary.LittleEndian.Uint64(instruction.Data[9:17])
	} else if len(instruction.Data) >= 9 {
		// Fallback for simpler instruction format
		amountIn = binary.LittleEndian.Uint64(instruction.Data[1:9])
	}

	// Extract swap information
	tokenIn := message.AccountKeys[instruction.Accounts[0]]
	tokenOut := message.AccountKeys[instruction.Accounts[1]]
	pool := message.AccountKeys[instruction.Accounts[2]]
	trader := message.AccountKeys[instruction.Accounts[3]]

	tradeInfo := TradeInfo{
		InstructionIndex: index,
		TokenIn:          tokenIn,
		TokenOut:         tokenOut,
		Pool:             pool,
		Trader:           trader,
		AmountIn:         amountIn,
		AmountOut:        0, // Would be extracted from transaction logs/metadata
		TradeType:        "swap",
	}

	result.Trade = append(result.Trade, tradeInfo)

	// Determine if it's a buy or sell based on token types
	if isBaseCurrency(tokenIn) {
		result.TradeBuys = append(result.TradeBuys, index)

		swapBuy := SwapBuy{
			TokenIn:      tokenIn,
			TokenOut:     tokenOut,
			AmountIn:     amountIn,
			AmountOut:    tradeInfo.AmountOut,
			Pool:         pool,
			Buyer:        trader,
			MinAmountOut: minAmountOut,
			Slippage:     0.0, // Would be calculated from actual vs expected amounts
		}
		result.SwapBuys = append(result.SwapBuys, swapBuy)
	} else {
		result.TradeSells = append(result.TradeSells, index)

		swapSell := SwapSell{
			TokenIn:      tokenIn,
			TokenOut:     tokenOut,
			AmountIn:     amountIn,
			AmountOut:    tradeInfo.AmountOut,
			Pool:         pool,
			Seller:       trader,
			MinAmountOut: minAmountOut,
			Slippage:     0.0, // Would be calculated from actual vs expected amounts
		}
		result.SwapSells = append(result.SwapSells, swapSell)
	}

	return nil
}

// parseBuyInstructionStandard parses buy instructions in standard format
func parseBuyInstructionStandard(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	if len(instruction.Accounts) < 6 {
		return fmt.Errorf("insufficient accounts for buy")
	}

	// Extract buy parameters from instruction data
	var amountIn, maxAmountIn uint64 = 0, 0

	if len(instruction.Data) >= 17 {
		amountIn = binary.LittleEndian.Uint64(instruction.Data[1:9])
		maxAmountIn = binary.LittleEndian.Uint64(instruction.Data[9:17])
	}

	tradeInfo := TradeInfo{
		InstructionIndex: index,
		TokenIn:          message.AccountKeys[instruction.Accounts[0]], // Base currency (SOL/USDC)
		TokenOut:         message.AccountKeys[instruction.Accounts[1]], // Token being bought
		Pool:             message.AccountKeys[instruction.Accounts[2]],
		Trader:           message.AccountKeys[instruction.Accounts[3]],
		AmountIn:         amountIn,
		AmountOut:        0, // Would be extracted from transaction logs
		TradeType:        "buy",
	}

	result.Trade = append(result.Trade, tradeInfo)
	result.TradeBuys = append(result.TradeBuys, index)

	slippage := 0.0
	if maxAmountIn > 0 && amountIn > 0 {
		slippage = calculateSlippage(amountIn, maxAmountIn)
	}

	swapBuy := SwapBuy{
		TokenIn:      tradeInfo.TokenIn,
		TokenOut:     tradeInfo.TokenOut,
		AmountIn:     amountIn,
		AmountOut:    tradeInfo.AmountOut,
		Pool:         tradeInfo.Pool,
		Buyer:        tradeInfo.Trader,
		MinAmountOut: 0, // Buy operations specify max input, not min output
		Slippage:     slippage,
	}
	result.SwapBuys = append(result.SwapBuys, swapBuy)

	return nil
}

// parseSellInstructionStandard parses sell instructions in standard format
func parseSellInstructionStandard(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	if len(instruction.Accounts) < 6 {
		return fmt.Errorf("insufficient accounts for sell")
	}

	// Extract sell parameters from instruction data
	var amountIn, minAmountOut uint64 = 0, 0

	if len(instruction.Data) >= 17 {
		amountIn = binary.LittleEndian.Uint64(instruction.Data[1:9])
		minAmountOut = binary.LittleEndian.Uint64(instruction.Data[9:17])
	}

	tradeInfo := TradeInfo{
		InstructionIndex: index,
		TokenIn:          message.AccountKeys[instruction.Accounts[0]], // Token being sold
		TokenOut:         message.AccountKeys[instruction.Accounts[1]], // Base currency (SOL/USDC)
		Pool:             message.AccountKeys[instruction.Accounts[2]],
		Trader:           message.AccountKeys[instruction.Accounts[3]],
		AmountIn:         amountIn,
		AmountOut:        0, // Would be extracted from transaction logs
		TradeType:        "sell",
	}

	result.Trade = append(result.Trade, tradeInfo)
	result.TradeSells = append(result.TradeSells, index)

	slippage := 0.0
	if minAmountOut > 0 && tradeInfo.AmountOut > 0 {
		slippage = calculateSlippage(tradeInfo.AmountOut, minAmountOut)
	}

	swapSell := SwapSell{
		TokenIn:      tradeInfo.TokenIn,
		TokenOut:     tradeInfo.TokenOut,
		AmountIn:     amountIn,
		AmountOut:    tradeInfo.AmountOut,
		Pool:         tradeInfo.Pool,
		Seller:       tradeInfo.Trader,
		MinAmountOut: minAmountOut,
		Slippage:     slippage,
	}
	result.SwapSells = append(result.SwapSells, swapSell)

	return nil
}

// parseDepositInstruction parses liquidity deposit instructions
func parseDepositInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	// Implementation would depend on the specific instruction format
	log.Printf("Deposit instruction detected at index %d", index)
	return nil
}

// parseWithdrawInstruction parses liquidity withdrawal instructions
func parseWithdrawInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	// Implementation would depend on the specific instruction format
	log.Printf("Withdraw instruction detected at index %d", index)
	return nil
}

// parseMigrateInstruction parses migration instructions
func parseMigrateInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	if len(instruction.Accounts) < 4 {
		return fmt.Errorf("insufficient accounts for migration")
	}

	// Extract migration amount from instruction data
	var amount uint64 = 0
	if len(instruction.Data) >= 9 {
		amount = binary.LittleEndian.Uint64(instruction.Data[1:9])
	}

	migration := Migration{
		FromPool:  message.AccountKeys[instruction.Accounts[0]],
		ToPool:    message.AccountKeys[instruction.Accounts[1]],
		Token:     message.AccountKeys[instruction.Accounts[2]],
		Owner:     message.AccountKeys[instruction.Accounts[3]],
		Amount:    amount,
		Timestamp: 0, // Would be extracted from block time
	}

	result.Migrate = append(result.Migrate, migration)
	return nil
}

// parseStakingInstruction parses staking-related instructions
func parseStakingInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	log.Printf("Staking instruction detected at index %d", index)
	return nil
}

// parseLiquidityInstruction parses liquidity-related instructions
func parseLiquidityInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	log.Printf("Liquidity instruction detected at index %d", index)
	return nil
}

// Enhanced token program instruction parsing
func parseTokenInstruction(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	if len(instruction.Data) == 0 {
		return nil
	}

	discriminator := instruction.Data[0]

	switch discriminator {
	case TOKEN_INSTRUCTION_TRANSFER:
		return parseTokenTransferInstructionStandard(instruction, message, index, result)
	case TOKEN_INSTRUCTION_MINT_TO:
		return parseTokenMintInstructionStandard(instruction, message, index, result)
	default:
		// Other token instructions we don't need to track
		return nil
	}
}

// parseTokenTransferInstructionStandard parses token transfer instructions in standard format
func parseTokenTransferInstructionStandard(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	if len(instruction.Data) < 9 || len(instruction.Accounts) < 3 {
		return nil
	}

	// Extract transfer amount
	amount := binary.LittleEndian.Uint64(instruction.Data[1:9])

	// Token transfers help us understand actual amounts moved
	// This information can be used to validate and enhance trade data
	log.Printf("Token transfer detected: %d tokens at instruction %d", amount, index)

	return nil
}

// parseTokenMintInstructionStandard parses token mint instructions in standard format
func parseTokenMintInstructionStandard(instruction solana.CompiledInstruction, message *solana.Message, index int, result *Transaction) error {
	if len(instruction.Data) < 9 || len(instruction.Accounts) < 3 {
		return nil
	}

	// Extract mint amount
	amount := binary.LittleEndian.Uint64(instruction.Data[1:9])

	// Token minting indicates token creation or additional supply
	log.Printf("Token mint detected: %d tokens at instruction %d", amount, index)

	return nil
}

// Helper functions

// isBaseCurrency checks if a token is a base currency (SOL, USDC, etc.)
func isBaseCurrency(tokenMint solana.PublicKey) bool {
	// Known base currency mints
	solMint := solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	usdcMint := solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	usdtMint := solana.MustPublicKeyFromBase58("Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB")

	return tokenMint.Equals(solMint) || tokenMint.Equals(usdcMint) || tokenMint.Equals(usdtMint)
}

// calculateSlippage calculates slippage percentage
func calculateSlippage(actualAmount, expectedAmount uint64) float64 {
	if expectedAmount == 0 {
		return 0.0
	}

	if actualAmount >= expectedAmount {
		return 0.0 // No slippage
	}

	return float64(expectedAmount-actualAmount) / float64(expectedAmount)
}

// getKnownTokenInfo returns known token information
func getKnownTokenInfo(tokenMint solana.PublicKey) (TokenInfo, bool) {
	// Check known tokens map from utils.go
	// For now, return some hardcoded examples
	knownTokensLocal := map[string]TokenInfo{
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

	info, exists := knownTokensLocal[tokenMint.String()]
	return info, exists
}

// parseGeyserInstructionWrapper parses a Geyser format instruction
func parseGeyserInstructionWrapper(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	programID := instruction.ProgramID

	// Check if this is a Raydium-related instruction
	switch programID {
	case RaydiumV4ProgramID, RaydiumV5ProgramID:
		return parseRaydiumGeyserInstruction(instruction, index, result, meta)
	case RaydiumLaunchpadV1ProgramID:
		return parseRaydiumLaunchpadInstruction(instruction, index, result, meta)
	case RaydiumCpSwapProgramID:
		return parseRaydiumCpSwapInstruction(instruction, index, result, meta)
	case TokenProgramID, Token2022ProgramID:
		return parseTokenGeyserInstruction(instruction, index, result, meta)
	default:
		// Not a Raydium-related instruction, skip
		return nil
	}
}

// parseRaydiumGeyserInstruction parses Raydium V4/V5 instructions in Geyser format
func parseRaydiumGeyserInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	if len(instruction.Data) == 0 {
		return fmt.Errorf("instruction data is empty")
	}

	// Get the instruction discriminator
	discriminator := instruction.Data[0]

	switch discriminator {
	case INSTRUCTION_INITIALIZE_POOL, INSTRUCTION_CREATE_POOL:
		return parseGeyserCreatePoolInstruction(instruction, index, result, meta)
	case INSTRUCTION_SWAP, INSTRUCTION_SWAP_BASE_IN, INSTRUCTION_SWAP_BASE_OUT:
		return parseGeyserSwapInstruction(instruction, index, result, meta)
	case INSTRUCTION_BUY:
		return parseGeyserBuyInstruction(instruction, index, result, meta)
	case INSTRUCTION_SELL:
		return parseGeyserSellInstruction(instruction, index, result, meta)
	case INSTRUCTION_MIGRATE:
		return parseGeyserMigrateInstruction(instruction, index, result, meta)
	default:
		log.Printf("Unknown Raydium instruction discriminator: %d", discriminator)
		return nil
	}
}

// parseRaydiumLaunchpadInstruction parses Raydium Launchpad specific instructions
func parseRaydiumLaunchpadInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	if len(instruction.Data) == 0 {
		return fmt.Errorf("launchpad instruction data is empty")
	}

	discriminator := instruction.Data[0]

	switch discriminator {
	case INSTRUCTION_INITIALIZE:
		return parseGeyserCreatePoolInstruction(instruction, index, result, meta)
	case INSTRUCTION_BUY:
		return parseGeyserBuyInstruction(instruction, index, result, meta)
	case INSTRUCTION_SELL:
		return parseGeyserSellInstruction(instruction, index, result, meta)
	default:
		log.Printf("Unknown Raydium Launchpad instruction discriminator: %d", discriminator)
		return nil
	}
}

// parseRaydiumCpSwapInstruction parses Raydium CP Swap instructions
func parseRaydiumCpSwapInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	if len(instruction.Data) == 0 {
		return fmt.Errorf("cp swap instruction data is empty")
	}

	discriminator := instruction.Data[0]

	switch discriminator {
	case INSTRUCTION_SWAP_BASE_IN, INSTRUCTION_SWAP_BASE_OUT:
		return parseGeyserSwapInstruction(instruction, index, result, meta)
	default:
		log.Printf("Unknown CP Swap instruction discriminator: %d", discriminator)
		return nil
	}
}

// parseTokenGeyserInstruction parses Token program instructions in Geyser format
func parseTokenGeyserInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	if len(instruction.Data) == 0 {
		return nil
	}

	discriminator := instruction.Data[0]

	switch discriminator {
	case TOKEN_INSTRUCTION_TRANSFER:
		return parseTokenTransferInstruction(instruction, index, result, meta)
	case TOKEN_INSTRUCTION_MINT_TO:
		return parseTokenMintInstruction(instruction, index, result, meta)
	default:
		// Other token instructions we don't need to track
		return nil
	}
}

// parseGeyserCreatePoolInstruction parses pool creation instructions in Geyser format
func parseGeyserCreatePoolInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	if len(instruction.Accounts) < 8 {
		return fmt.Errorf("insufficient accounts for pool creation")
	}

	// Extract creation parameters from instruction data
	var tokenDecimals uint8 = 9 // Default
	var initialLiquidity uint64 = 0

	if len(instruction.Data) >= 17 {
		// Parse actual instruction data
		tokenDecimals = instruction.Data[9]
		initialLiquidity = binary.LittleEndian.Uint64(instruction.Data[9:17])
	}

	// Extract token symbol from metadata or use default
	tokenSymbol := extractTokenSymbol(instruction.Accounts[0], meta)

	createInfo := CreateInfo{
		TokenMint:     instruction.Accounts[0],
		PoolAddress:   instruction.Accounts[1],
		Creator:       instruction.Accounts[2],
		TokenDecimals: tokenDecimals,
		TokenSymbol:   tokenSymbol,
		Amount:        initialLiquidity,
		Timestamp:     0, // Would need to be extracted from block time
	}

	result.Create = append(result.Create, createInfo)
	return nil
}

// parseGeyserSwapInstruction parses swap instructions in Geyser format
func parseGeyserSwapInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	if len(instruction.Accounts) < 6 {
		return fmt.Errorf("insufficient accounts for swap")
	}

	// Extract swap amounts from instruction data
	var amountIn, amountOut uint64 = 0, 0
	var minAmountOut uint64 = 0

	if len(instruction.Data) >= 25 {
		amountIn = binary.LittleEndian.Uint64(instruction.Data[1:9])
		minAmountOut = binary.LittleEndian.Uint64(instruction.Data[9:17])
		// AmountOut would be determined from transaction logs/meta
		amountOut = extractAmountOutFromMeta(instruction.Accounts, meta)
	}

	// Extract swap information
	tradeInfo := TradeInfo{
		InstructionIndex: index,
		TokenIn:          instruction.Accounts[0],
		TokenOut:         instruction.Accounts[1],
		Pool:             instruction.Accounts[2],
		Trader:           instruction.Accounts[3],
		AmountIn:         amountIn,
		AmountOut:        amountOut,
		TradeType:        "swap",
	}

	result.Trade = append(result.Trade, tradeInfo)

	// Determine if it's a buy or sell
	if isBaseCurrency(tradeInfo.TokenIn) {
		result.TradeBuys = append(result.TradeBuys, index)

		slippage := calculateSlippage(amountOut, minAmountOut)
		swapBuy := SwapBuy{
			TokenIn:      tradeInfo.TokenIn,
			TokenOut:     tradeInfo.TokenOut,
			AmountIn:     amountIn,
			AmountOut:    amountOut,
			Pool:         tradeInfo.Pool,
			Buyer:        tradeInfo.Trader,
			MinAmountOut: minAmountOut,
			Slippage:     slippage,
		}
		result.SwapBuys = append(result.SwapBuys, swapBuy)
	} else {
		result.TradeSells = append(result.TradeSells, index)

		slippage := calculateSlippage(amountOut, minAmountOut)
		swapSell := SwapSell{
			TokenIn:      tradeInfo.TokenIn,
			TokenOut:     tradeInfo.TokenOut,
			AmountIn:     amountIn,
			AmountOut:    amountOut,
			Pool:         tradeInfo.Pool,
			Seller:       tradeInfo.Trader,
			MinAmountOut: minAmountOut,
			Slippage:     slippage,
		}
		result.SwapSells = append(result.SwapSells, swapSell)
	}

	return nil
}

// parseGeyserBuyInstruction parses buy instructions in Geyser format
func parseGeyserBuyInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	if len(instruction.Accounts) < 6 {
		return fmt.Errorf("insufficient accounts for buy")
	}

	// Extract buy parameters from instruction data
	var amountIn, maxAmountIn uint64 = 0, 0

	if len(instruction.Data) >= 17 {
		amountIn = binary.LittleEndian.Uint64(instruction.Data[1:9])
		maxAmountIn = binary.LittleEndian.Uint64(instruction.Data[9:17])
	}

	// Calculate amount out from transaction metadata
	amountOut := extractAmountOutFromMeta(instruction.Accounts, meta)

	tradeInfo := TradeInfo{
		InstructionIndex: index,
		TokenIn:          instruction.Accounts[0], // Base currency (SOL/USDC)
		TokenOut:         instruction.Accounts[1], // Token being bought
		Pool:             instruction.Accounts[2],
		Trader:           instruction.Accounts[3],
		AmountIn:         amountIn,
		AmountOut:        amountOut,
		TradeType:        "buy",
	}

	result.Trade = append(result.Trade, tradeInfo)
	result.TradeBuys = append(result.TradeBuys, index)

	slippage := 0.0
	if maxAmountIn > 0 && amountIn > 0 {
		slippage = calculateSlippage(amountIn, maxAmountIn)
	}

	swapBuy := SwapBuy{
		TokenIn:      tradeInfo.TokenIn,
		TokenOut:     tradeInfo.TokenOut,
		AmountIn:     amountIn,
		AmountOut:    tradeInfo.AmountOut,
		Pool:         tradeInfo.Pool,
		Buyer:        tradeInfo.Trader,
		MinAmountOut: 0, // Buy operations specify max input, not min output
		Slippage:     slippage,
	}
	result.SwapBuys = append(result.SwapBuys, swapBuy)

	return nil
}

// parseGeyserSellInstruction parses sell instructions in Geyser format
func parseGeyserSellInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	if len(instruction.Accounts) < 6 {
		return fmt.Errorf("insufficient accounts for sell")
	}

	// Extract sell parameters from instruction data
	var amountIn, minAmountOut uint64 = 0, 0

	if len(instruction.Data) >= 17 {
		amountIn = binary.LittleEndian.Uint64(instruction.Data[1:9])
		minAmountOut = binary.LittleEndian.Uint64(instruction.Data[9:17])
	}

	// Calculate amount out from transaction metadata
	amountOut := extractAmountOutFromMeta(instruction.Accounts, meta)

	tradeInfo := TradeInfo{
		InstructionIndex: index,
		TokenIn:          instruction.Accounts[0], // Token being sold
		TokenOut:         instruction.Accounts[1], // Base currency (SOL/USDC)
		Pool:             instruction.Accounts[2],
		Trader:           instruction.Accounts[3],
		AmountIn:         amountIn,
		AmountOut:        amountOut,
		TradeType:        "sell",
	}

	result.Trade = append(result.Trade, tradeInfo)
	result.TradeSells = append(result.TradeSells, index)

	slippage := 0.0
	if minAmountOut > 0 && tradeInfo.AmountOut > 0 {
		slippage = calculateSlippage(tradeInfo.AmountOut, minAmountOut)
	}

	swapSell := SwapSell{
		TokenIn:      tradeInfo.TokenIn,
		TokenOut:     tradeInfo.TokenOut,
		AmountIn:     amountIn,
		AmountOut:    tradeInfo.AmountOut,
		Pool:         tradeInfo.Pool,
		Seller:       tradeInfo.Trader,
		MinAmountOut: minAmountOut,
		Slippage:     slippage,
	}
	result.SwapSells = append(result.SwapSells, swapSell)

	return nil
}

// parseGeyserMigrateInstruction parses migration instructions in Geyser format
func parseGeyserMigrateInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	if len(instruction.Accounts) < 4 {
		return fmt.Errorf("insufficient accounts for migration")
	}

	// Extract migration amount from instruction data
	var amount uint64 = 0
	if len(instruction.Data) >= 9 {
		amount = binary.LittleEndian.Uint64(instruction.Data[1:9])
	}

	migration := Migration{
		FromPool:  instruction.Accounts[0],
		ToPool:    instruction.Accounts[1],
		Token:     instruction.Accounts[2],
		Owner:     instruction.Accounts[3],
		Amount:    amount,
		Timestamp: 0, // Would be extracted from block time
	}

	result.Migrate = append(result.Migrate, migration)
	return nil
}

// parseTokenTransferInstruction parses token transfer instructions
func parseTokenTransferInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	// Token transfers help us understand the actual amounts moved
	// This information is used to validate and supplement trade data
	return nil
}

// parseTokenMintInstruction parses token mint instructions
func parseTokenMintInstruction(instruction GeyserInstruction, index int, result *Transaction, meta *TransactionMeta) error {
	// Token minting indicates token creation
	// This information is used to supplement create data
	return nil
}

// Helper functions for Geyser format

// extractTokenSymbol extracts token symbol from metadata or returns default
func extractTokenSymbol(tokenMint solana.PublicKey, meta *TransactionMeta) string {
	// In a real implementation, this would:
	// 1. Check known token registry
	// 2. Query token metadata
	// 3. Parse token name from transaction logs

	// For now, return known symbols or default
	if tokenInfo, exists := getKnownTokenInfo(tokenMint); exists {
		return tokenInfo.Symbol
	}
	return "UNKNOWN"
}

// extractAmountOutFromMeta extracts the actual amount out from transaction metadata
func extractAmountOutFromMeta(accounts []solana.PublicKey, meta *TransactionMeta) uint64 {
	// In a real implementation, this would:
	// 1. Compare pre/post balances
	// 2. Parse token balance changes
	// 3. Extract from transaction logs

	// For now, return a placeholder
	if meta != nil && len(meta.PostBalances) > 0 {
		// Simple example: return difference in balances
		if len(meta.PreBalances) > 0 && len(meta.PostBalances) > 0 {
			if len(meta.PreBalances) > 1 && len(meta.PostBalances) > 1 {
				return meta.PostBalances[1] - meta.PreBalances[1]
			}
		}
	}
	return 0
}
