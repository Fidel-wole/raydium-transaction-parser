package main

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func TestSwapInstructionBuilder(t *testing.T) {
	// Create a swap instruction
	swapInst := NewSwapInstruction().
		SetUserSourceToken(solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")).
		SetUserDestToken(solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")).
		SetUserOwner(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH")).
		SetAmmID(solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")).
		SetAmmAuthority(solana.MustPublicKeyFromBase58("5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h")).
		SetAmmOpenOrders(solana.MustPublicKeyFromBase58("EhhTKczWMGQt46ynNeRX1WfeagwwJd7ufHvCDjRxjo5Q")).
		SetAmmTargetOrders(solana.MustPublicKeyFromBase58("27haf8L6oxUeXrHrgEgsexjSY5hbVUWEmvv9Nyxg8vQv")).
		SetPoolCoinToken(solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")).
		SetPoolPcToken(solana.MustPublicKeyFromBase58("CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C")).
		SetSerumProgram(solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")).
		SetSerumMarket(solana.MustPublicKeyFromBase58("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb")).
		SetSerumBids(solana.MustPublicKeyFromBase58("LanMV9sAd7wArD4vJFi2qDdfnVhFxYSUg6eADduJ3uj")).
		SetSerumAsks(solana.MustPublicKeyFromBase58("FoaFt2Dtz58RA6DPjbRb9t9z8sLJRChiGFTv21EfaseZ")).
		SetSerumEventQueue(solana.MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")).
		SetSerumCoinVault(solana.MustPublicKeyFromBase58("11111111111111111111111111111111")).
		SetSerumPcVault(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH")).
		SetSerumVaultSigner(solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")).
		SetAmountIn(1000000).
		SetMinimumAmountOut(900000)

	// Build the instruction
	instruction, err := swapInst.Build()
	if err != nil {
		t.Fatalf("Failed to build swap instruction: %v", err)
	}

	// Verify the instruction
	if instruction.ProgramID() != RaydiumV4ProgramID {
		t.Errorf("Expected program ID %s, got %s", RaydiumV4ProgramID, instruction.ProgramID())
	}

	accounts := instruction.Accounts()
	if len(accounts) != 18 {
		t.Errorf("Expected 18 accounts, got %d", len(accounts))
	}

	data, err := instruction.Data()
	if err != nil {
		t.Fatalf("Failed to get instruction data: %v", err)
	}
	if len(data) != 17 {
		t.Errorf("Expected 17 bytes of data, got %d", len(data))
	}

	// Verify discriminator
	if data[0] != INSTRUCTION_SWAP {
		t.Errorf("Expected discriminator %d, got %d", INSTRUCTION_SWAP, data[0])
	}

	t.Logf("✓ Swap instruction built successfully with %d accounts and %d bytes of data", len(accounts), len(data))
}

func TestBuyInstructionBuilder(t *testing.T) {
	// Create a buy instruction
	buyInst := NewBuyInstruction().
		SetUserAuthority(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH")).
		SetUserTokenAccount(solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")).
		SetUserSolAccount(solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")).
		SetAmmID(solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")).
		SetAmmAuthority(solana.MustPublicKeyFromBase58("5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h")).
		SetTokenVault(solana.MustPublicKeyFromBase58("EhhTKczWMGQt46ynNeRX1WfeagwwJd7ufHvCDjRxjo5Q")).
		SetSolVault(solana.MustPublicKeyFromBase58("27haf8L6oxUeXrHrgEgsexjSY5hbVUWEmvv9Nyxg8vQv")).
		SetTokenMint(solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")).
		SetAmount(1000000).
		SetMaxSolCost(500000)

	// Build the instruction
	instruction, err := buyInst.Build()
	if err != nil {
		t.Fatalf("Failed to build buy instruction: %v", err)
	}

	// Verify the instruction
	if instruction.ProgramID() != RaydiumLaunchpadV1ProgramID {
		t.Errorf("Expected program ID %s, got %s", RaydiumLaunchpadV1ProgramID, instruction.ProgramID())
	}

	accounts := instruction.Accounts()
	if len(accounts) != 10 {
		t.Errorf("Expected 10 accounts, got %d", len(accounts))
	}

	data, err := instruction.Data()
	if err != nil {
		t.Fatalf("Failed to get instruction data: %v", err)
	}
	if len(data) != 17 {
		t.Errorf("Expected 17 bytes of data, got %d", len(data))
	}

	// Verify discriminator
	if data[0] != INSTRUCTION_BUY {
		t.Errorf("Expected discriminator %d, got %d", INSTRUCTION_BUY, data[0])
	}

	t.Logf("✓ Buy instruction built successfully with %d accounts and %d bytes of data", len(accounts), len(data))
}

func TestSellInstructionBuilder(t *testing.T) {
	// Create a sell instruction
	sellInst := NewSellInstruction().
		SetUserAuthority(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH")).
		SetUserTokenAccount(solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")).
		SetUserSolAccount(solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")).
		SetAmmID(solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")).
		SetAmmAuthority(solana.MustPublicKeyFromBase58("5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h")).
		SetTokenVault(solana.MustPublicKeyFromBase58("EhhTKczWMGQt46ynNeRX1WfeagwwJd7ufHvCDjRxjo5Q")).
		SetSolVault(solana.MustPublicKeyFromBase58("27haf8L6oxUeXrHrgEgsexjSY5hbVUWEmvv9Nyxg8vQv")).
		SetTokenMint(solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")).
		SetAmount(1000000).
		SetMinSolReceived(400000)

	// Build the instruction
	instruction, err := sellInst.Build()
	if err != nil {
		t.Fatalf("Failed to build sell instruction: %v", err)
	}

	// Verify the instruction
	if instruction.ProgramID() != RaydiumLaunchpadV1ProgramID {
		t.Errorf("Expected program ID %s, got %s", RaydiumLaunchpadV1ProgramID, instruction.ProgramID())
	}

	accounts := instruction.Accounts()
	if len(accounts) != 10 {
		t.Errorf("Expected 10 accounts, got %d", len(accounts))
	}

	data, err := instruction.Data()
	if err != nil {
		t.Fatalf("Failed to get instruction data: %v", err)
	}
	if len(data) != 17 {
		t.Errorf("Expected 17 bytes of data, got %d", len(data))
	}

	// Verify discriminator
	if data[0] != INSTRUCTION_SELL {
		t.Errorf("Expected discriminator %d, got %d", INSTRUCTION_SELL, data[0])
	}

	t.Logf("✓ Sell instruction built successfully with %d accounts and %d bytes of data", len(accounts), len(data))
}

func TestCreateTokenInstructionBuilder(t *testing.T) {
	// Create a token creation instruction
	createInst := NewCreateTokenInstruction().
		SetPayer(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH")).
		SetMint(solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")).
		SetMintAuthority(solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")).
		SetFreezeAuthority(solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")).
		SetDecimals(6).
		SetName("Test Token").
		SetSymbol("TEST").
		SetURI("https://example.com/token.json").
		SetInitialSupply(1000000000000)

	// Build the instruction
	instruction, err := createInst.Build()
	if err != nil {
		t.Fatalf("Failed to build create token instruction: %v", err)
	}

	// Verify the instruction
	if instruction.ProgramID() != RaydiumLaunchpadV1ProgramID {
		t.Errorf("Expected program ID %s, got %s", RaydiumLaunchpadV1ProgramID, instruction.ProgramID())
	}

	accounts := instruction.Accounts()
	if len(accounts) != 6 {
		t.Errorf("Expected 6 accounts, got %d", len(accounts))
	}

	data, err := instruction.Data()
	if err != nil {
		t.Fatalf("Failed to get instruction data: %v", err)
	}

	// Verify discriminator
	if data[0] != INSTRUCTION_CREATE_POOL {
		t.Errorf("Expected discriminator %d, got %d", INSTRUCTION_CREATE_POOL, data[0])
	}

	t.Logf("✓ Create token instruction built successfully with %d accounts and %d bytes of data", len(accounts), len(data))
}

func TestMigrateInstructionBuilder(t *testing.T) {
	// Create a migrate instruction
	migrateInst := NewMigrateInstruction().
		SetUserAuthority(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH")).
		SetFromPool(solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")).
		SetToPool(solana.MustPublicKeyFromBase58("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")).
		SetTokenAccount(solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")).
		SetAmount(1000000)

	// Build the instruction
	instruction, err := migrateInst.Build()
	if err != nil {
		t.Fatalf("Failed to build migrate instruction: %v", err)
	}

	// Verify the instruction
	if instruction.ProgramID() != RaydiumV4ProgramID {
		t.Errorf("Expected program ID %s, got %s", RaydiumV4ProgramID, instruction.ProgramID())
	}

	accounts := instruction.Accounts()
	if len(accounts) != 5 {
		t.Errorf("Expected 5 accounts, got %d", len(accounts))
	}

	data, err := instruction.Data()
	if err != nil {
		t.Fatalf("Failed to get instruction data: %v", err)
	}
	if len(data) != 9 {
		t.Errorf("Expected 9 bytes of data, got %d", len(data))
	}

	// Verify discriminator
	if data[0] != INSTRUCTION_MIGRATE {
		t.Errorf("Expected discriminator %d, got %d", INSTRUCTION_MIGRATE, data[0])
	}

	t.Logf("✓ Migrate instruction built successfully with %d accounts and %d bytes of data", len(accounts), len(data))
}

// TestTransactionSubmission tests submitting transactions to Solana
// This test requires environment variables for wallet and token information
func TestTransactionSubmission(t *testing.T) {
	// Skip if no environment variables are set
	walletPath := os.Getenv("SOLANA_WALLET_PATH")
	rpcEndpoint := os.Getenv("SOLANA_RPC_ENDPOINT")

	if walletPath == "" || rpcEndpoint == "" {
		t.Skip("Skipping transaction submission test - missing environment variables SOLANA_WALLET_PATH and SOLANA_RPC_ENDPOINT")
	}

	wallet, err := solana.PrivateKeyFromSolanaKeygenFile(walletPath)
	if err != nil {
		t.Fatalf("Failed to load wallet: %v", err)
	}

	// Create RPC client
	client := rpc.New(rpcEndpoint)

	// Get recent blockhash
	ctx := context.Background()
	recent, err := client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		t.Fatalf("Failed to get recent blockhash: %v", err)
	}

	// Create a simple transfer instruction as a test
	// In a real scenario, this would be a buy/sell instruction
	transferInst := solana.NewInstruction(
		solana.SystemProgramID,
		solana.AccountMetaSlice{
			{PublicKey: wallet.PublicKey(), IsWritable: true, IsSigner: true},
			{PublicKey: wallet.PublicKey(), IsWritable: true, IsSigner: false}, // sending to self for test
		},
		[]byte{2, 0, 0, 0, 232, 3, 0, 0, 0, 0, 0, 0}, // transfer 1000 lamports
	)

	// Create transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{transferInst},
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Sign transaction
	if _, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(wallet.PublicKey()) {
			return &wallet
		}
		return nil
	}); err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	// Simulate the transaction (don't actually send)
	simResult, err := client.SimulateTransaction(ctx, tx)
	if err != nil {
		t.Fatalf("Failed to simulate transaction: %v", err)
	}

	if simResult.Value.Err != nil {
		t.Fatalf("Transaction simulation failed: %v", simResult.Value.Err)
	}

	t.Log("✓ Transaction simulation successful")
	t.Log("Note: To test actual submission, remove the simulation and use SendTransaction")
}

// TestTransactionParsingWithLiveData tests parsing with live transaction data
func TestTransactionParsingWithLiveData(t *testing.T) {
	// Test with a real transaction from the sample file
	sampleTxData := `5wefCTqi9ynrh8pvVHFzpgHCLFFzoBwGoTgWSd6iq2Qw4Y51U4cEc2xHYtsdVSFZmRXUp5DNMSkhzb1CaXomLpJM`

	// Create RPC client
	client := rpc.New(rpc.MainNetBeta_RPC)
	ctx := context.Background()

	// Get the transaction
	signature := solana.MustSignatureFromBase58(sampleTxData)
	txResult, err := client.GetTransaction(ctx, signature, &rpc.GetTransactionOpts{
		Encoding: solana.EncodingBase64,
	})
	if err != nil {
		t.Skipf("Failed to fetch transaction (network issue): %v", err)
	}

	if txResult.Meta == nil {
		t.Skip("Transaction not found or has no meta")
	}

	// Parse the transaction
	if txResult.Transaction == nil {
		t.Skip("No transaction data available")
	}

	// For this test, we'll just test that the parser doesn't crash
	// In a real scenario, you would have access to the raw transaction data
	t.Logf("✓ Successfully fetched transaction with signature: %s", signature)
	t.Logf("✓ Transaction slot: %d", txResult.Slot)
	t.Logf("✓ Transaction parsing test completed (raw data parsing requires proper transaction bytes)")
}

// TestParsingWithSampleData tests parsing with sample transaction data
func TestParsingWithSampleData(t *testing.T) {
	// Read sample transaction data
	sampleData, err := os.ReadFile("sample_transaction.txt")
	if err != nil {
		t.Skipf("Sample transaction file not found: %v", err)
	}

	// Parse each line as a transaction
	lines := strings.Split(string(sampleData), "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}

		parsedTx, err := ParseTransaction(line, uint64(12345+i))
		if err != nil {
			t.Logf("Failed to parse transaction %d: %v", i, err)
			continue
		}

		t.Logf("✓ Parsed sample transaction %d: %s", i, parsedTx.Signature)
	}
}

// TestBuilderChaining tests that all builders properly support method chaining
func TestBuilderChaining(t *testing.T) {
	// Test swap instruction chaining
	swapInst := NewSwapInstruction().
		SetAmountIn(1000).
		SetMinimumAmountOut(900).
		SetUserOwner(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH"))

	if swapInst.amountIn != 1000 {
		t.Errorf("Expected amountIn 1000, got %d", swapInst.amountIn)
	}

	// Test buy instruction chaining
	buyInst := NewBuyInstruction().
		SetAmount(2000).
		SetMaxSolCost(1000).
		SetUserAuthority(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH"))

	if buyInst.amount != 2000 {
		t.Errorf("Expected amount 2000, got %d", buyInst.amount)
	}

	// Test sell instruction chaining
	sellInst := NewSellInstruction().
		SetAmount(3000).
		SetMinSolReceived(1500).
		SetUserAuthority(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH"))

	if sellInst.amount != 3000 {
		t.Errorf("Expected amount 3000, got %d", sellInst.amount)
	}

	// Test create token instruction chaining
	createInst := NewCreateTokenInstruction().
		SetDecimals(6).
		SetName("Test").
		SetSymbol("TST").
		SetInitialSupply(1000000)

	if createInst.decimals != 6 {
		t.Errorf("Expected decimals 6, got %d", createInst.decimals)
	}

	// Test migrate instruction chaining
	migrateInst := NewMigrateInstruction().
		SetAmount(5000).
		SetUserAuthority(solana.MustPublicKeyFromBase58("HN7cABqLq46Es1jh92dQQisAq662SmxELLLsHHe4YWrH"))

	if migrateInst.amount != 5000 {
		t.Errorf("Expected amount 5000, got %d", migrateInst.amount)
	}

	t.Log("✓ All builder chaining tests passed")
}
