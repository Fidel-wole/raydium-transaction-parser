package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
)

// InstructionDebugInfo contains comprehensive debugging information for each instruction
type InstructionDebugInfo struct {
	InstructionIndex int                   `json:"instruction_index"`
	ProgramID        string                `json:"program_id"`
	Discriminator    string                `json:"discriminator"`
	DataLength       int                   `json:"data_length"`
	AccountCount     int                   `json:"account_count"`
	Accounts         []AccountDebugInfo    `json:"accounts"`
	Parameters       InstructionParameters `json:"parameters"`
	ParsedResult     interface{}           `json:"parsed_result"`
}

// AccountDebugInfo contains detailed information about each account
type AccountDebugInfo struct {
	Index       int    `json:"index"`
	Address     string `json:"address"`
	Description string `json:"description"`
	IsSystem    bool   `json:"is_system"`
	IsProgram   bool   `json:"is_program"`
	IsToken     bool   `json:"is_token"`
	IsSigner    bool   `json:"is_signer"`
}

// InstructionParameters contains parsed parameters from instruction data
type InstructionParameters struct {
	RawData       []byte                 `json:"raw_data"`
	Discriminator uint64                 `json:"discriminator"`
	Amount        uint64                 `json:"amount"`
	MaxAmount     uint64                 `json:"max_amount"`
	MinAmount     uint64                 `json:"min_amount"`
	Decimals      uint8                  `json:"decimals"`
	Timestamp     int64                  `json:"timestamp"`
	ExtraParams   map[string]interface{} `json:"extra_params"`
}

// TransactionDebugInfo contains comprehensive debugging information for the entire transaction
type TransactionDebugInfo struct {
	Signature    string                 `json:"signature"`
	Slot         uint64                 `json:"slot"`
	Timestamp    int64                  `json:"timestamp"`
	AllAccounts  []AccountDebugInfo     `json:"all_accounts"`
	Instructions []InstructionDebugInfo `json:"instructions"`
	Summary      TransactionSummary     `json:"summary"`
}

// TransactionSummary provides a high-level summary of the transaction
type TransactionSummary struct {
	TotalInstructions   int `json:"total_instructions"`
	RaydiumInstructions int `json:"raydium_instructions"`
	TokenProgram        int `json:"token_program"`
	SystemProgram       int `json:"system_program"`
	CreateOps           int `json:"create_ops"`
	TradeOps            int `json:"trade_ops"`
	SwapOps             int `json:"swap_ops"`
	MigrateOps          int `json:"migrate_ops"`
}

// Enhanced token info with more details
type EnhancedTokenInfo struct {
	Mint        string `json:"mint"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Decimals    uint8  `json:"decimals"`
	Supply      uint64 `json:"supply"`
	IsKnown     bool   `json:"is_known"`
	Description string `json:"description"`
}

// Function to get enhanced token information
func getEnhancedTokenInfo(tokenMint solana.PublicKey) EnhancedTokenInfo {
	knownTokens := map[string]EnhancedTokenInfo{
		"So11111111111111111111111111111111111111112": {
			Mint:        "So11111111111111111111111111111111111111112",
			Symbol:      "SOL",
			Name:        "Solana",
			Decimals:    9,
			Supply:      0, // Dynamic supply
			IsKnown:     true,
			Description: "Native Solana token",
		},
		"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": {
			Mint:        "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			Symbol:      "USDC",
			Name:        "USD Coin",
			Decimals:    6,
			Supply:      0, // Dynamic supply
			IsKnown:     true,
			Description: "USD Coin stablecoin",
		},
		"8pf71rxkus6HVhNa9ERdJ571wfPa1a8QKKMsxGkDbonk": {
			Mint:        "8pf71rxkus6HVhNa9ERdJ571wfPa1a8QKKMsxGkDbonk",
			Symbol:      "JAMAL", // Based on the transaction context
			Name:        "Jamal Token",
			Decimals:    6, // Common default, should be parsed from metadata
			Supply:      0,
			IsKnown:     true,
			Description: "Raydium Launchpad token",
		},
	}

	if info, exists := knownTokens[tokenMint.String()]; exists {
		return info
	}

	// Return unknown token info
	return EnhancedTokenInfo{
		Mint:        tokenMint.String(),
		Symbol:      "UNKNOWN",
		Name:        "Unknown Token",
		Decimals:    6, // Default decimals
		Supply:      0,
		IsKnown:     false,
		Description: "Unknown token",
	}
}

// Function to classify account type
func classifyAccount(account solana.PublicKey) AccountDebugInfo {
	address := account.String()

	info := AccountDebugInfo{
		Address:     address,
		Description: "Unknown account",
		IsSystem:    false,
		IsProgram:   false,
		IsToken:     false,
		IsSigner:    false,
	}

	// Check if it's a system account
	if account.Equals(SystemProgramID) {
		info.Description = "System Program"
		info.IsSystem = true
		info.IsProgram = true
	} else if account.Equals(TokenProgramID) {
		info.Description = "Token Program"
		info.IsSystem = true
		info.IsProgram = true
	} else if account.Equals(AssociatedTokenProgramID) {
		info.Description = "Associated Token Program"
		info.IsSystem = true
		info.IsProgram = true
	} else if account.Equals(RaydiumLaunchpadV1ProgramID) {
		info.Description = "Raydium Launchpad V1 Program"
		info.IsProgram = true
	} else if account.Equals(RaydiumCpSwapProgramID) {
		info.Description = "Raydium CP Swap Program"
		info.IsProgram = true
	} else if account.Equals(RaydiumV4ProgramID) {
		info.Description = "Raydium V4 Program"
		info.IsProgram = true
	} else if account.Equals(RaydiumV5ProgramID) {
		info.Description = "Raydium V5 Program"
		info.IsProgram = true
	} else if account.String() == "So11111111111111111111111111111111111111112" {
		info.Description = "SOL (Wrapped SOL)"
		info.IsToken = true
	} else if account.String() == "8pf71rxkus6HVhNa9ERdJ571wfPa1a8QKKMsxGkDbonk" {
		info.Description = "JAMAL Token Mint"
		info.IsToken = true
	} else {
		// Try to determine if it's a token account or pool
		if len(address) == 44 { // Standard Solana address length
			info.Description = "Token Account / Pool / User Account"
		}
	}

	return info
}

// Function to create comprehensive debug info for a transaction
func createTransactionDebugInfo(tx *Transaction, message *solana.Message) *TransactionDebugInfo {
	debugInfo := &TransactionDebugInfo{
		Signature:    tx.Signature.String(),
		Slot:         tx.Slot,
		Timestamp:    time.Now().Unix(),
		AllAccounts:  make([]AccountDebugInfo, 0),
		Instructions: make([]InstructionDebugInfo, 0),
		Summary: TransactionSummary{
			CreateOps:  len(tx.Create),
			TradeOps:   len(tx.Trade),
			SwapOps:    len(tx.SwapBuys) + len(tx.SwapSells),
			MigrateOps: len(tx.Migrate),
		},
	}

	// Process all accounts
	for i, account := range message.AccountKeys {
		accountInfo := classifyAccount(account)
		accountInfo.Index = i
		accountInfo.IsSigner = (i == 0) // First account is typically the signer
		debugInfo.AllAccounts = append(debugInfo.AllAccounts, accountInfo)
	}

	return debugInfo
}

// Function to add instruction debug info
func addInstructionDebugInfo(debugInfo *TransactionDebugInfo, instruction solana.CompiledInstruction, message *solana.Message, index int, programID solana.PublicKey) {
	instrInfo := InstructionDebugInfo{
		InstructionIndex: index,
		ProgramID:        programID.String(),
		DataLength:       len(instruction.Data),
		AccountCount:     len(instruction.Accounts),
		Accounts:         make([]AccountDebugInfo, 0),
		Parameters: InstructionParameters{
			RawData:     instruction.Data,
			ExtraParams: make(map[string]interface{}),
		},
	}

	// Parse discriminator
	if len(instruction.Data) >= 8 {
		instrInfo.Discriminator = fmt.Sprintf("%x", instruction.Data[:8])
		instrInfo.Parameters.Discriminator = binary.LittleEndian.Uint64(instruction.Data[:8])
	} else if len(instruction.Data) >= 1 {
		instrInfo.Discriminator = fmt.Sprintf("%x", instruction.Data[0])
		instrInfo.Parameters.Discriminator = uint64(instruction.Data[0])
	}

	// Parse amounts if present
	if len(instruction.Data) >= 16 {
		instrInfo.Parameters.Amount = binary.LittleEndian.Uint64(instruction.Data[8:16])
	}
	if len(instruction.Data) >= 24 {
		instrInfo.Parameters.MaxAmount = binary.LittleEndian.Uint64(instruction.Data[16:24])
	}
	if len(instruction.Data) >= 17 {
		instrInfo.Parameters.Decimals = instruction.Data[16]
	}

	// Process instruction accounts
	for i, accountIndex := range instruction.Accounts {
		if int(accountIndex) < len(message.AccountKeys) {
			accountInfo := classifyAccount(message.AccountKeys[accountIndex])
			accountInfo.Index = i
			instrInfo.Accounts = append(instrInfo.Accounts, accountInfo)
		}
	}

	debugInfo.Instructions = append(debugInfo.Instructions, instrInfo)
	debugInfo.Summary.TotalInstructions++

	// Update summary based on program type
	if programID.Equals(RaydiumLaunchpadV1ProgramID) || programID.Equals(RaydiumV4ProgramID) || programID.Equals(RaydiumV5ProgramID) {
		debugInfo.Summary.RaydiumInstructions++
	} else if programID.Equals(TokenProgramID) {
		debugInfo.Summary.TokenProgram++
	} else if programID.Equals(SystemProgramID) {
		debugInfo.Summary.SystemProgram++
	}
}

// Function to print comprehensive debug info
func printTransactionDebugInfo(debugInfo *TransactionDebugInfo) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("COMPREHENSIVE TRANSACTION DEBUG INFO")
	fmt.Println(strings.Repeat("=", 80))

	jsonData, err := json.MarshalIndent(debugInfo, "", "  ")
	if err != nil {
		log.Printf("Error marshaling debug info to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
	fmt.Println(strings.Repeat("=", 80))
}

// Enhanced debug structure for comprehensive instruction debugging
type ComprehensiveInstructionDebug struct {
	InstructionIndex int                     `json:"instruction_index"`
	ProgramID        string                  `json:"program_id"`
	ProgramName      string                  `json:"program_name"`
	Discriminator    string                  `json:"discriminator"`
	DataLength       int                     `json:"data_length"`
	DataHex          string                  `json:"data_hex"`
	AccountCount     int                     `json:"account_count"`
	Accounts         []DetailedAccountInfo   `json:"accounts"`
	Parameters       ComprehensiveParameters `json:"parameters"`
	ParsedResult     interface{}             `json:"parsed_result"`
	Timestamp        int64                   `json:"timestamp"`
}

// Enhanced account info with all 18 fields
type DetailedAccountInfo struct {
	Index         int    `json:"index"`
	Address       string `json:"address"`
	Description   string `json:"description"`
	Role          string `json:"role"`
	IsSystem      bool   `json:"is_system"`
	IsProgram     bool   `json:"is_program"`
	IsToken       bool   `json:"is_token"`
	IsSigner      bool   `json:"is_signer"`
	IsWritable    bool   `json:"is_writable"`
	IsExecutable  bool   `json:"is_executable"`
	IsOwner       bool   `json:"is_owner"`
	IsRentExempt  bool   `json:"is_rent_exempt"`
	Balance       uint64 `json:"balance"`
	DataSize      uint64 `json:"data_size"`
	TokenMint     string `json:"token_mint"`
	TokenOwner    string `json:"token_owner"`
	TokenAmount   uint64 `json:"token_amount"`
	TokenDecimals uint8  `json:"token_decimals"`
}

// Comprehensive parameters structure
type ComprehensiveParameters struct {
	RawData       []byte                 `json:"raw_data"`
	DataHex       string                 `json:"data_hex"`
	Discriminator uint64                 `json:"discriminator"`
	Amount        uint64                 `json:"amount"`
	AmountOut     uint64                 `json:"amount_out"`
	MaxAmount     uint64                 `json:"max_amount"`
	MinAmount     uint64                 `json:"min_amount"`
	Decimals      uint8                  `json:"decimals"`
	Timestamp     int64                  `json:"timestamp"`
	Slippage      float64                `json:"slippage"`
	Direction     string                 `json:"direction"`
	TokenIn       string                 `json:"token_in"`
	TokenOut      string                 `json:"token_out"`
	Pool          string                 `json:"pool"`
	Creator       string                 `json:"creator"`
	Trader        string                 `json:"trader"`
	ExtraParams   map[string]interface{} `json:"extra_params"`
}

// Function to create comprehensive instruction debug info
func createInstructionDebugInfo(instruction solana.CompiledInstruction, message *solana.Message, index int, programID solana.PublicKey) *ComprehensiveInstructionDebug {
	debugInfo := &ComprehensiveInstructionDebug{
		InstructionIndex: index,
		ProgramID:        programID.String(),
		ProgramName:      getProgramName(programID),
		DataLength:       len(instruction.Data),
		DataHex:          fmt.Sprintf("%x", instruction.Data),
		AccountCount:     len(instruction.Accounts),
		Accounts:         make([]DetailedAccountInfo, 0),
		Parameters: ComprehensiveParameters{
			RawData:     instruction.Data,
			DataHex:     fmt.Sprintf("%x", instruction.Data),
			ExtraParams: make(map[string]interface{}),
		},
		Timestamp: time.Now().Unix(),
	}

	// Parse discriminator
	if len(instruction.Data) >= 8 {
		debugInfo.Discriminator = fmt.Sprintf("%x", instruction.Data[:8])
		debugInfo.Parameters.Discriminator = binary.LittleEndian.Uint64(instruction.Data[:8])
	} else if len(instruction.Data) >= 1 {
		debugInfo.Discriminator = fmt.Sprintf("%x", instruction.Data[0])
		debugInfo.Parameters.Discriminator = uint64(instruction.Data[0])
	}

	// Parse amounts based on program type
	if programID.Equals(RaydiumLaunchpadV1ProgramID) {
		parseRaydiumLaunchpadParameters(&debugInfo.Parameters, instruction.Data)
	} else if programID.Equals(RaydiumV4ProgramID) || programID.Equals(RaydiumV5ProgramID) {
		parseRaydiumV4V5Parameters(&debugInfo.Parameters, instruction.Data)
	} else if programID.Equals(TokenProgramID) {
		parseTokenProgramParameters(&debugInfo.Parameters, instruction.Data)
	}

	// Process all accounts with comprehensive info
	for i, accountIndex := range instruction.Accounts {
		if int(accountIndex) < len(message.AccountKeys) {
			account := message.AccountKeys[accountIndex]
			accountInfo := createDetailedAccountInfo(account, i, int(accountIndex), programID)
			debugInfo.Accounts = append(debugInfo.Accounts, accountInfo)
		}
	}

	return debugInfo
}

// Function to get program name
func getProgramName(programID solana.PublicKey) string {
	switch programID {
	case RaydiumV4ProgramID:
		return "Raydium V4"
	case RaydiumV5ProgramID:
		return "Raydium V5"
	case RaydiumLaunchpadV1ProgramID:
		return "Raydium Launchpad V1"
	case RaydiumCpSwapProgramID:
		return "Raydium CP Swap"
	case RaydiumStakingProgramID:
		return "Raydium Staking"
	case RaydiumLiquidityProgramID:
		return "Raydium Liquidity"
	case TokenProgramID:
		return "Token Program"
	case SystemProgramID:
		return "System Program"
	case AssociatedTokenProgramID:
		return "Associated Token Program"
	default:
		return "Unknown Program"
	}
}

// Function to create detailed account info with all 18 fields
func createDetailedAccountInfo(account solana.PublicKey, instructionIndex int, accountIndex int, programID solana.PublicKey) DetailedAccountInfo {
	address := account.String()

	info := DetailedAccountInfo{
		Index:         instructionIndex,
		Address:       address,
		Description:   "Unknown account",
		Role:          "unknown",
		IsSystem:      false,
		IsProgram:     false,
		IsToken:       false,
		IsSigner:      (accountIndex == 0), // First account is typically the signer
		IsWritable:    false,
		IsExecutable:  false,
		IsOwner:       false,
		IsRentExempt:  false,
		Balance:       0,
		DataSize:      0,
		TokenMint:     "",
		TokenOwner:    "",
		TokenAmount:   0,
		TokenDecimals: 0,
	}

	// Classify account type and set appropriate fields
	if account.Equals(SystemProgramID) {
		info.Description = "System Program"
		info.Role = "system_program"
		info.IsSystem = true
		info.IsProgram = true
		info.IsExecutable = true
	} else if account.Equals(TokenProgramID) {
		info.Description = "Token Program"
		info.Role = "token_program"
		info.IsSystem = true
		info.IsProgram = true
		info.IsExecutable = true
	} else if account.Equals(AssociatedTokenProgramID) {
		info.Description = "Associated Token Program"
		info.Role = "ata_program"
		info.IsSystem = true
		info.IsProgram = true
		info.IsExecutable = true
	} else if account.Equals(RaydiumLaunchpadV1ProgramID) {
		info.Description = "Raydium Launchpad V1 Program"
		info.Role = "launchpad_program"
		info.IsProgram = true
		info.IsExecutable = true
	} else if account.Equals(RaydiumCpSwapProgramID) {
		info.Description = "Raydium CP Swap Program"
		info.Role = "cpswap_program"
		info.IsProgram = true
		info.IsExecutable = true
	} else if account.Equals(RaydiumV4ProgramID) {
		info.Description = "Raydium V4 Program"
		info.Role = "raydium_v4_program"
		info.IsProgram = true
		info.IsExecutable = true
	} else if account.Equals(RaydiumV5ProgramID) {
		info.Description = "Raydium V5 Program"
		info.Role = "raydium_v5_program"
		info.IsProgram = true
		info.IsExecutable = true
	} else if account.String() == "So11111111111111111111111111111111111111112" {
		info.Description = "SOL (Wrapped SOL)"
		info.Role = "sol_mint"
		info.IsToken = true
		info.TokenMint = address
		info.TokenDecimals = 9
	} else if account.String() == "8pf71rxkus6HVhNa9ERdJ571wfPa1a8QKKMsxGkDbonk" {
		info.Description = "JAMAL Token Mint"
		info.Role = "token_mint"
		info.IsToken = true
		info.TokenMint = address
		info.TokenDecimals = 6
	} else {
		// Try to determine role based on context
		if programID.Equals(RaydiumLaunchpadV1ProgramID) {
			info.Role = determineRaydiumLaunchpadRole(instructionIndex, account)
		} else if programID.Equals(RaydiumV4ProgramID) || programID.Equals(RaydiumV5ProgramID) {
			info.Role = determineRaydiumV4V5Role(instructionIndex, account)
		} else {
			info.Role = "user_account"
		}
		info.Description = fmt.Sprintf("Account (%s)", info.Role)
		info.IsWritable = true // Most user accounts are writable
	}

	return info
}

// Function to determine role in Raydium Launchpad instructions
func determineRaydiumLaunchpadRole(accountIndex int, account solana.PublicKey) string {
	switch accountIndex {
	case 0:
		return "signer/creator"
	case 1:
		return "token_mint"
	case 2:
		return "pool"
	case 3:
		return "user_token_account"
	case 4:
		return "user_sol_account"
	case 5:
		return "pool_sol_account"
	case 6:
		return "pool_token_account"
	case 7:
		return "system_program"
	case 8:
		return "token_program"
	case 9:
		return "associated_token_program"
	default:
		return "additional_account"
	}
}

// Function to determine role in Raydium V4/V5 instructions
func determineRaydiumV4V5Role(accountIndex int, account solana.PublicKey) string {
	switch accountIndex {
	case 0:
		return "signer/trader"
	case 1:
		return "user_source_account"
	case 2:
		return "user_dest_account"
	case 3:
		return "pool_source_account"
	case 4:
		return "pool_dest_account"
	case 5:
		return "pool_authority"
	case 6:
		return "amm_pool"
	case 7:
		return "amm_authority"
	case 8:
		return "amm_open_orders"
	case 9:
		return "amm_target_orders"
	default:
		return "additional_account"
	}
}

// Function to parse Raydium Launchpad parameters
func parseRaydiumLaunchpadParameters(params *ComprehensiveParameters, data []byte) {
	// For Launchpad transactions, we need to skip the discriminator
	dataStart := 1
	if len(data) >= 8 {
		dataStart = 8 // Skip 8-byte discriminator for complex instructions
	}

	if len(data) >= dataStart+8 {
		params.Amount = binary.LittleEndian.Uint64(data[dataStart : dataStart+8])
	}
	if len(data) >= dataStart+16 {
		params.MaxAmount = binary.LittleEndian.Uint64(data[dataStart+8 : dataStart+16])
	}
	if len(data) >= dataStart+17 {
		params.Decimals = data[dataStart+16]
	}
	params.Direction = "buy" // Default for Launchpad
}

// Function to parse Raydium V4/V5 parameters
func parseRaydiumV4V5Parameters(params *ComprehensiveParameters, data []byte) {
	if len(data) >= 9 {
		params.Amount = binary.LittleEndian.Uint64(data[1:9])
	}
	if len(data) >= 17 {
		params.MinAmount = binary.LittleEndian.Uint64(data[9:17])
	}
	if len(data) >= 18 {
		params.Decimals = data[17]
	}
}

// Function to parse Token Program parameters
func parseTokenProgramParameters(params *ComprehensiveParameters, data []byte) {
	if len(data) >= 9 {
		params.Amount = binary.LittleEndian.Uint64(data[1:9])
	}
	if len(data) >= 17 {
		params.AmountOut = binary.LittleEndian.Uint64(data[9:17])
	}
}

// Function to print comprehensive instruction debug info
func printInstructionDebugInfo(debugInfo *ComprehensiveInstructionDebug) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 100))
	fmt.Printf("INSTRUCTION DEBUG INFO - Index: %d | Program: %s\n", debugInfo.InstructionIndex, debugInfo.ProgramName)
	fmt.Printf("%s\n", strings.Repeat("=", 100))

	// Print basic instruction info
	fmt.Printf("Program ID: %s\n", debugInfo.ProgramID)
	fmt.Printf("Discriminator: %s\n", debugInfo.Discriminator)
	fmt.Printf("Data Length: %d bytes\n", debugInfo.DataLength)
	fmt.Printf("Data Hex: %s\n", debugInfo.DataHex)
	fmt.Printf("Account Count: %d\n", debugInfo.AccountCount)
	fmt.Printf("Timestamp: %d\n", debugInfo.Timestamp)

	// Print all account details
	fmt.Printf("\nALL ACCOUNT DETAILS (18 fields per account):\n")
	fmt.Printf("%s\n", strings.Repeat("-", 100))
	for i, account := range debugInfo.Accounts {
		fmt.Printf("Account %d:\n", i)
		fmt.Printf("  Address: %s\n", account.Address)
		fmt.Printf("  Description: %s\n", account.Description)
		fmt.Printf("  Role: %s\n", account.Role)
		fmt.Printf("  IsSystem: %t\n", account.IsSystem)
		fmt.Printf("  IsProgram: %t\n", account.IsProgram)
		fmt.Printf("  IsToken: %t\n", account.IsToken)
		fmt.Printf("  IsSigner: %t\n", account.IsSigner)
		fmt.Printf("  IsWritable: %t\n", account.IsWritable)
		fmt.Printf("  IsExecutable: %t\n", account.IsExecutable)
		fmt.Printf("  IsOwner: %t\n", account.IsOwner)
		fmt.Printf("  IsRentExempt: %t\n", account.IsRentExempt)
		fmt.Printf("  Balance: %d\n", account.Balance)
		fmt.Printf("  DataSize: %d\n", account.DataSize)
		fmt.Printf("  TokenMint: %s\n", account.TokenMint)
		fmt.Printf("  TokenOwner: %s\n", account.TokenOwner)
		fmt.Printf("  TokenAmount: %d\n", account.TokenAmount)
		fmt.Printf("  TokenDecimals: %d\n", account.TokenDecimals)
		fmt.Printf("%s\n", strings.Repeat("-", 50))
	}

	// Print parsed parameters
	fmt.Printf("\nPARSED PARAMETERS:\n")
	fmt.Printf("Amount: %d\n", debugInfo.Parameters.Amount)
	fmt.Printf("AmountOut: %d\n", debugInfo.Parameters.AmountOut)
	fmt.Printf("MaxAmount: %d\n", debugInfo.Parameters.MaxAmount)
	fmt.Printf("MinAmount: %d\n", debugInfo.Parameters.MinAmount)
	fmt.Printf("Decimals: %d\n", debugInfo.Parameters.Decimals)
	fmt.Printf("Direction: %s\n", debugInfo.Parameters.Direction)
	fmt.Printf("TokenIn: %s\n", debugInfo.Parameters.TokenIn)
	fmt.Printf("TokenOut: %s\n", debugInfo.Parameters.TokenOut)
	fmt.Printf("Pool: %s\n", debugInfo.Parameters.Pool)
	fmt.Printf("Creator: %s\n", debugInfo.Parameters.Creator)
	fmt.Printf("Trader: %s\n", debugInfo.Parameters.Trader)
	fmt.Printf("Slippage: %.4f\n", debugInfo.Parameters.Slippage)

	// Print full JSON
	fmt.Printf("\nFULL JSON OUTPUT:\n")
	jsonData, err := json.MarshalIndent(debugInfo, "", "  ")
	if err != nil {
		log.Printf("Error marshaling debug info to JSON: %v", err)
	} else {
		fmt.Println(string(jsonData))
	}

	fmt.Printf("\n%s\n", strings.Repeat("=", 100))
}
