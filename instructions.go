package main

import (
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

// SwapInstruction represents a Raydium swap instruction
type SwapInstruction struct {
	programID        solana.PublicKey
	userSourceToken  solana.PublicKey
	userDestToken    solana.PublicKey
	userOwner        solana.PublicKey
	ammID            solana.PublicKey
	ammAuthority     solana.PublicKey
	ammOpenOrders    solana.PublicKey
	ammTargetOrders  solana.PublicKey
	poolCoinToken    solana.PublicKey
	poolPcToken      solana.PublicKey
	serumProgram     solana.PublicKey
	serumMarket      solana.PublicKey
	serumBids        solana.PublicKey
	serumAsks        solana.PublicKey
	serumEventQueue  solana.PublicKey
	serumCoinVault   solana.PublicKey
	serumPcVault     solana.PublicKey
	serumVaultSigner solana.PublicKey
	amountIn         uint64
	minimumAmountOut uint64
}

// NewSwapInstruction creates a new swap instruction builder
func NewSwapInstruction() *SwapInstruction {
	return &SwapInstruction{
		programID: RaydiumV4ProgramID,
	}
}

// SetProgramID sets the program ID for the swap instruction
func (s *SwapInstruction) SetProgramID(programID solana.PublicKey) *SwapInstruction {
	s.programID = programID
	return s
}

// SetUserSourceToken sets the user's source token account
func (s *SwapInstruction) SetUserSourceToken(userSourceToken solana.PublicKey) *SwapInstruction {
	s.userSourceToken = userSourceToken
	return s
}

// SetUserDestToken sets the user's destination token account
func (s *SwapInstruction) SetUserDestToken(userDestToken solana.PublicKey) *SwapInstruction {
	s.userDestToken = userDestToken
	return s
}

// SetUserOwner sets the user/owner account
func (s *SwapInstruction) SetUserOwner(userOwner solana.PublicKey) *SwapInstruction {
	s.userOwner = userOwner
	return s
}

// SetAmmID sets the AMM ID
func (s *SwapInstruction) SetAmmID(ammID solana.PublicKey) *SwapInstruction {
	s.ammID = ammID
	return s
}

// SetAmmAuthority sets the AMM authority
func (s *SwapInstruction) SetAmmAuthority(ammAuthority solana.PublicKey) *SwapInstruction {
	s.ammAuthority = ammAuthority
	return s
}

// SetAmmOpenOrders sets the AMM open orders account
func (s *SwapInstruction) SetAmmOpenOrders(ammOpenOrders solana.PublicKey) *SwapInstruction {
	s.ammOpenOrders = ammOpenOrders
	return s
}

// SetAmmTargetOrders sets the AMM target orders account
func (s *SwapInstruction) SetAmmTargetOrders(ammTargetOrders solana.PublicKey) *SwapInstruction {
	s.ammTargetOrders = ammTargetOrders
	return s
}

// SetPoolCoinToken sets the pool coin token account
func (s *SwapInstruction) SetPoolCoinToken(poolCoinToken solana.PublicKey) *SwapInstruction {
	s.poolCoinToken = poolCoinToken
	return s
}

// SetPoolPcToken sets the pool PC token account
func (s *SwapInstruction) SetPoolPcToken(poolPcToken solana.PublicKey) *SwapInstruction {
	s.poolPcToken = poolPcToken
	return s
}

// SetSerumProgram sets the Serum program ID
func (s *SwapInstruction) SetSerumProgram(serumProgram solana.PublicKey) *SwapInstruction {
	s.serumProgram = serumProgram
	return s
}

// SetSerumMarket sets the Serum market account
func (s *SwapInstruction) SetSerumMarket(serumMarket solana.PublicKey) *SwapInstruction {
	s.serumMarket = serumMarket
	return s
}

// SetSerumBids sets the Serum bids account
func (s *SwapInstruction) SetSerumBids(serumBids solana.PublicKey) *SwapInstruction {
	s.serumBids = serumBids
	return s
}

// SetSerumAsks sets the Serum asks account
func (s *SwapInstruction) SetSerumAsks(serumAsks solana.PublicKey) *SwapInstruction {
	s.serumAsks = serumAsks
	return s
}

// SetSerumEventQueue sets the Serum event queue account
func (s *SwapInstruction) SetSerumEventQueue(serumEventQueue solana.PublicKey) *SwapInstruction {
	s.serumEventQueue = serumEventQueue
	return s
}

// SetSerumCoinVault sets the Serum coin vault account
func (s *SwapInstruction) SetSerumCoinVault(serumCoinVault solana.PublicKey) *SwapInstruction {
	s.serumCoinVault = serumCoinVault
	return s
}

// SetSerumPcVault sets the Serum PC vault account
func (s *SwapInstruction) SetSerumPcVault(serumPcVault solana.PublicKey) *SwapInstruction {
	s.serumPcVault = serumPcVault
	return s
}

// SetSerumVaultSigner sets the Serum vault signer
func (s *SwapInstruction) SetSerumVaultSigner(serumVaultSigner solana.PublicKey) *SwapInstruction {
	s.serumVaultSigner = serumVaultSigner
	return s
}

// SetAmountIn sets the amount to swap in
func (s *SwapInstruction) SetAmountIn(amountIn uint64) *SwapInstruction {
	s.amountIn = amountIn
	return s
}

// SetMinimumAmountOut sets the minimum amount out
func (s *SwapInstruction) SetMinimumAmountOut(minimumAmountOut uint64) *SwapInstruction {
	s.minimumAmountOut = minimumAmountOut
	return s
}

// Build creates the Solana instruction
func (s *SwapInstruction) Build() (solana.Instruction, error) {
	// Build instruction data
	data := make([]byte, 17) // 1 byte discriminator + 8 bytes amountIn + 8 bytes minimumAmountOut
	data[0] = INSTRUCTION_SWAP
	binary.LittleEndian.PutUint64(data[1:9], s.amountIn)
	binary.LittleEndian.PutUint64(data[9:17], s.minimumAmountOut)

	// Build accounts slice
	accounts := solana.AccountMetaSlice{
		{PublicKey: s.userSourceToken, IsWritable: true, IsSigner: false},
		{PublicKey: s.userDestToken, IsWritable: true, IsSigner: false},
		{PublicKey: s.userOwner, IsWritable: false, IsSigner: true},
		{PublicKey: s.ammID, IsWritable: true, IsSigner: false},
		{PublicKey: s.ammAuthority, IsWritable: false, IsSigner: false},
		{PublicKey: s.ammOpenOrders, IsWritable: true, IsSigner: false},
		{PublicKey: s.ammTargetOrders, IsWritable: true, IsSigner: false},
		{PublicKey: s.poolCoinToken, IsWritable: true, IsSigner: false},
		{PublicKey: s.poolPcToken, IsWritable: true, IsSigner: false},
		{PublicKey: s.serumProgram, IsWritable: false, IsSigner: false},
		{PublicKey: s.serumMarket, IsWritable: true, IsSigner: false},
		{PublicKey: s.serumBids, IsWritable: true, IsSigner: false},
		{PublicKey: s.serumAsks, IsWritable: true, IsSigner: false},
		{PublicKey: s.serumEventQueue, IsWritable: true, IsSigner: false},
		{PublicKey: s.serumCoinVault, IsWritable: true, IsSigner: false},
		{PublicKey: s.serumPcVault, IsWritable: true, IsSigner: false},
		{PublicKey: s.serumVaultSigner, IsWritable: false, IsSigner: false},
		{PublicKey: TokenProgramID, IsWritable: false, IsSigner: false},
	}

	return solana.NewInstruction(
		s.programID,
		accounts,
		data,
	), nil
}

// BuyInstruction represents a Raydium buy instruction
type BuyInstruction struct {
	programID        solana.PublicKey
	userAuthority    solana.PublicKey
	userTokenAccount solana.PublicKey
	userSolAccount   solana.PublicKey
	ammID            solana.PublicKey
	ammAuthority     solana.PublicKey
	tokenVault       solana.PublicKey
	solVault         solana.PublicKey
	tokenMint        solana.PublicKey
	amount           uint64
	maxSolCost       uint64
}

// NewBuyInstruction creates a new buy instruction builder
func NewBuyInstruction() *BuyInstruction {
	return &BuyInstruction{
		programID: RaydiumLaunchpadV1ProgramID,
	}
}

// SetProgramID sets the program ID for the buy instruction
func (b *BuyInstruction) SetProgramID(programID solana.PublicKey) *BuyInstruction {
	b.programID = programID
	return b
}

// SetUserAuthority sets the user authority
func (b *BuyInstruction) SetUserAuthority(userAuthority solana.PublicKey) *BuyInstruction {
	b.userAuthority = userAuthority
	return b
}

// SetUserTokenAccount sets the user token account
func (b *BuyInstruction) SetUserTokenAccount(userTokenAccount solana.PublicKey) *BuyInstruction {
	b.userTokenAccount = userTokenAccount
	return b
}

// SetUserSolAccount sets the user SOL account
func (b *BuyInstruction) SetUserSolAccount(userSolAccount solana.PublicKey) *BuyInstruction {
	b.userSolAccount = userSolAccount
	return b
}

// SetAmmID sets the AMM ID
func (b *BuyInstruction) SetAmmID(ammID solana.PublicKey) *BuyInstruction {
	b.ammID = ammID
	return b
}

// SetAmmAuthority sets the AMM authority
func (b *BuyInstruction) SetAmmAuthority(ammAuthority solana.PublicKey) *BuyInstruction {
	b.ammAuthority = ammAuthority
	return b
}

// SetTokenVault sets the token vault
func (b *BuyInstruction) SetTokenVault(tokenVault solana.PublicKey) *BuyInstruction {
	b.tokenVault = tokenVault
	return b
}

// SetSolVault sets the SOL vault
func (b *BuyInstruction) SetSolVault(solVault solana.PublicKey) *BuyInstruction {
	b.solVault = solVault
	return b
}

// SetTokenMint sets the token mint
func (b *BuyInstruction) SetTokenMint(tokenMint solana.PublicKey) *BuyInstruction {
	b.tokenMint = tokenMint
	return b
}

// SetAmount sets the amount to buy
func (b *BuyInstruction) SetAmount(amount uint64) *BuyInstruction {
	b.amount = amount
	return b
}

// SetMaxSolCost sets the maximum SOL cost
func (b *BuyInstruction) SetMaxSolCost(maxSolCost uint64) *BuyInstruction {
	b.maxSolCost = maxSolCost
	return b
}

// Build creates the Solana instruction
func (b *BuyInstruction) Build() (solana.Instruction, error) {
	// Build instruction data
	data := make([]byte, 17) // 1 byte discriminator + 8 bytes amount + 8 bytes maxSolCost
	data[0] = INSTRUCTION_BUY
	binary.LittleEndian.PutUint64(data[1:9], b.amount)
	binary.LittleEndian.PutUint64(data[9:17], b.maxSolCost)

	// Build accounts slice
	accounts := solana.AccountMetaSlice{
		{PublicKey: b.userAuthority, IsWritable: false, IsSigner: true},
		{PublicKey: b.userTokenAccount, IsWritable: true, IsSigner: false},
		{PublicKey: b.userSolAccount, IsWritable: true, IsSigner: false},
		{PublicKey: b.ammID, IsWritable: true, IsSigner: false},
		{PublicKey: b.ammAuthority, IsWritable: false, IsSigner: false},
		{PublicKey: b.tokenVault, IsWritable: true, IsSigner: false},
		{PublicKey: b.solVault, IsWritable: true, IsSigner: false},
		{PublicKey: b.tokenMint, IsWritable: false, IsSigner: false},
		{PublicKey: TokenProgramID, IsWritable: false, IsSigner: false},
		{PublicKey: SystemProgramID, IsWritable: false, IsSigner: false},
	}

	return solana.NewInstruction(
		b.programID,
		accounts,
		data,
	), nil
}

// SellInstruction represents a Raydium sell instruction
type SellInstruction struct {
	programID        solana.PublicKey
	userAuthority    solana.PublicKey
	userTokenAccount solana.PublicKey
	userSolAccount   solana.PublicKey
	ammID            solana.PublicKey
	ammAuthority     solana.PublicKey
	tokenVault       solana.PublicKey
	solVault         solana.PublicKey
	tokenMint        solana.PublicKey
	amount           uint64
	minSolReceived   uint64
}

// NewSellInstruction creates a new sell instruction builder
func NewSellInstruction() *SellInstruction {
	return &SellInstruction{
		programID: RaydiumLaunchpadV1ProgramID,
	}
}

// SetProgramID sets the program ID for the sell instruction
func (s *SellInstruction) SetProgramID(programID solana.PublicKey) *SellInstruction {
	s.programID = programID
	return s
}

// SetUserAuthority sets the user authority
func (s *SellInstruction) SetUserAuthority(userAuthority solana.PublicKey) *SellInstruction {
	s.userAuthority = userAuthority
	return s
}

// SetUserTokenAccount sets the user token account
func (s *SellInstruction) SetUserTokenAccount(userTokenAccount solana.PublicKey) *SellInstruction {
	s.userTokenAccount = userTokenAccount
	return s
}

// SetUserSolAccount sets the user SOL account
func (s *SellInstruction) SetUserSolAccount(userSolAccount solana.PublicKey) *SellInstruction {
	s.userSolAccount = userSolAccount
	return s
}

// SetAmmID sets the AMM ID
func (s *SellInstruction) SetAmmID(ammID solana.PublicKey) *SellInstruction {
	s.ammID = ammID
	return s
}

// SetAmmAuthority sets the AMM authority
func (s *SellInstruction) SetAmmAuthority(ammAuthority solana.PublicKey) *SellInstruction {
	s.ammAuthority = ammAuthority
	return s
}

// SetTokenVault sets the token vault
func (s *SellInstruction) SetTokenVault(tokenVault solana.PublicKey) *SellInstruction {
	s.tokenVault = tokenVault
	return s
}

// SetSolVault sets the SOL vault
func (s *SellInstruction) SetSolVault(solVault solana.PublicKey) *SellInstruction {
	s.solVault = solVault
	return s
}

// SetTokenMint sets the token mint
func (s *SellInstruction) SetTokenMint(tokenMint solana.PublicKey) *SellInstruction {
	s.tokenMint = tokenMint
	return s
}

// SetAmount sets the amount to sell
func (s *SellInstruction) SetAmount(amount uint64) *SellInstruction {
	s.amount = amount
	return s
}

// SetMinSolReceived sets the minimum SOL received
func (s *SellInstruction) SetMinSolReceived(minSolReceived uint64) *SellInstruction {
	s.minSolReceived = minSolReceived
	return s
}

// Build creates the Solana instruction
func (s *SellInstruction) Build() (solana.Instruction, error) {
	// Build instruction data
	data := make([]byte, 17) // 1 byte discriminator + 8 bytes amount + 8 bytes minSolReceived
	data[0] = INSTRUCTION_SELL
	binary.LittleEndian.PutUint64(data[1:9], s.amount)
	binary.LittleEndian.PutUint64(data[9:17], s.minSolReceived)

	// Build accounts slice
	accounts := solana.AccountMetaSlice{
		{PublicKey: s.userAuthority, IsWritable: false, IsSigner: true},
		{PublicKey: s.userTokenAccount, IsWritable: true, IsSigner: false},
		{PublicKey: s.userSolAccount, IsWritable: true, IsSigner: false},
		{PublicKey: s.ammID, IsWritable: true, IsSigner: false},
		{PublicKey: s.ammAuthority, IsWritable: false, IsSigner: false},
		{PublicKey: s.tokenVault, IsWritable: true, IsSigner: false},
		{PublicKey: s.solVault, IsWritable: true, IsSigner: false},
		{PublicKey: s.tokenMint, IsWritable: false, IsSigner: false},
		{PublicKey: TokenProgramID, IsWritable: false, IsSigner: false},
		{PublicKey: SystemProgramID, IsWritable: false, IsSigner: false},
	}

	return solana.NewInstruction(
		s.programID,
		accounts,
		data,
	), nil
}

// CreateTokenInstruction represents a token creation instruction
type CreateTokenInstruction struct {
	programID       solana.PublicKey
	payer           solana.PublicKey
	mint            solana.PublicKey
	mintAuthority   solana.PublicKey
	freezeAuthority solana.PublicKey
	decimals        uint8
	name            string
	symbol          string
	uri             string
	initialSupply   uint64
}

// NewCreateTokenInstruction creates a new token creation instruction builder
func NewCreateTokenInstruction() *CreateTokenInstruction {
	return &CreateTokenInstruction{
		programID: RaydiumLaunchpadV1ProgramID,
		decimals:  9, // Default to 9 decimals
	}
}

// SetProgramID sets the program ID for the create token instruction
func (c *CreateTokenInstruction) SetProgramID(programID solana.PublicKey) *CreateTokenInstruction {
	c.programID = programID
	return c
}

// SetPayer sets the payer account
func (c *CreateTokenInstruction) SetPayer(payer solana.PublicKey) *CreateTokenInstruction {
	c.payer = payer
	return c
}

// SetMint sets the mint account
func (c *CreateTokenInstruction) SetMint(mint solana.PublicKey) *CreateTokenInstruction {
	c.mint = mint
	return c
}

// SetMintAuthority sets the mint authority
func (c *CreateTokenInstruction) SetMintAuthority(mintAuthority solana.PublicKey) *CreateTokenInstruction {
	c.mintAuthority = mintAuthority
	return c
}

// SetFreezeAuthority sets the freeze authority
func (c *CreateTokenInstruction) SetFreezeAuthority(freezeAuthority solana.PublicKey) *CreateTokenInstruction {
	c.freezeAuthority = freezeAuthority
	return c
}

// SetDecimals sets the token decimals
func (c *CreateTokenInstruction) SetDecimals(decimals uint8) *CreateTokenInstruction {
	c.decimals = decimals
	return c
}

// SetName sets the token name
func (c *CreateTokenInstruction) SetName(name string) *CreateTokenInstruction {
	c.name = name
	return c
}

// SetSymbol sets the token symbol
func (c *CreateTokenInstruction) SetSymbol(symbol string) *CreateTokenInstruction {
	c.symbol = symbol
	return c
}

// SetURI sets the token URI
func (c *CreateTokenInstruction) SetURI(uri string) *CreateTokenInstruction {
	c.uri = uri
	return c
}

// SetInitialSupply sets the initial supply
func (c *CreateTokenInstruction) SetInitialSupply(initialSupply uint64) *CreateTokenInstruction {
	c.initialSupply = initialSupply
	return c
}

// Build creates the Solana instruction
func (c *CreateTokenInstruction) Build() (solana.Instruction, error) {
	// Build instruction data
	nameBytes := []byte(c.name)
	symbolBytes := []byte(c.symbol)
	uriBytes := []byte(c.uri)

	// Calculate total data size
	dataSize := 1 + // discriminator
		1 + // decimals
		4 + len(nameBytes) + // name length + name
		4 + len(symbolBytes) + // symbol length + symbol
		4 + len(uriBytes) + // uri length + uri
		8 // initial supply

	data := make([]byte, dataSize)
	offset := 0

	// Discriminator
	data[offset] = INSTRUCTION_CREATE_POOL
	offset++

	// Decimals
	data[offset] = c.decimals
	offset++

	// Name
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(nameBytes)))
	offset += 4
	copy(data[offset:offset+len(nameBytes)], nameBytes)
	offset += len(nameBytes)

	// Symbol
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(symbolBytes)))
	offset += 4
	copy(data[offset:offset+len(symbolBytes)], symbolBytes)
	offset += len(symbolBytes)

	// URI
	binary.LittleEndian.PutUint32(data[offset:offset+4], uint32(len(uriBytes)))
	offset += 4
	copy(data[offset:offset+len(uriBytes)], uriBytes)
	offset += len(uriBytes)

	// Initial supply
	binary.LittleEndian.PutUint64(data[offset:offset+8], c.initialSupply)

	// Build accounts slice
	accounts := solana.AccountMetaSlice{
		{PublicKey: c.payer, IsWritable: true, IsSigner: true},
		{PublicKey: c.mint, IsWritable: true, IsSigner: false},
		{PublicKey: c.mintAuthority, IsWritable: false, IsSigner: false},
		{PublicKey: c.freezeAuthority, IsWritable: false, IsSigner: false},
		{PublicKey: TokenProgramID, IsWritable: false, IsSigner: false},
		{PublicKey: SystemProgramID, IsWritable: false, IsSigner: false},
	}

	return solana.NewInstruction(
		c.programID,
		accounts,
		data,
	), nil
}

// MigrateInstruction represents a migration instruction
type MigrateInstruction struct {
	programID     solana.PublicKey
	userAuthority solana.PublicKey
	fromPool      solana.PublicKey
	toPool        solana.PublicKey
	tokenAccount  solana.PublicKey
	amount        uint64
}

// NewMigrateInstruction creates a new migrate instruction builder
func NewMigrateInstruction() *MigrateInstruction {
	return &MigrateInstruction{
		programID: RaydiumV4ProgramID,
	}
}

// SetProgramID sets the program ID for the migrate instruction
func (m *MigrateInstruction) SetProgramID(programID solana.PublicKey) *MigrateInstruction {
	m.programID = programID
	return m
}

// SetUserAuthority sets the user authority
func (m *MigrateInstruction) SetUserAuthority(userAuthority solana.PublicKey) *MigrateInstruction {
	m.userAuthority = userAuthority
	return m
}

// SetFromPool sets the from pool
func (m *MigrateInstruction) SetFromPool(fromPool solana.PublicKey) *MigrateInstruction {
	m.fromPool = fromPool
	return m
}

// SetToPool sets the to pool
func (m *MigrateInstruction) SetToPool(toPool solana.PublicKey) *MigrateInstruction {
	m.toPool = toPool
	return m
}

// SetTokenAccount sets the token account
func (m *MigrateInstruction) SetTokenAccount(tokenAccount solana.PublicKey) *MigrateInstruction {
	m.tokenAccount = tokenAccount
	return m
}

// SetAmount sets the amount to migrate
func (m *MigrateInstruction) SetAmount(amount uint64) *MigrateInstruction {
	m.amount = amount
	return m
}

// Build creates the Solana instruction
func (m *MigrateInstruction) Build() (solana.Instruction, error) {
	// Build instruction data
	data := make([]byte, 9) // 1 byte discriminator + 8 bytes amount
	data[0] = INSTRUCTION_MIGRATE
	binary.LittleEndian.PutUint64(data[1:9], m.amount)

	// Build accounts slice
	accounts := solana.AccountMetaSlice{
		{PublicKey: m.userAuthority, IsWritable: false, IsSigner: true},
		{PublicKey: m.fromPool, IsWritable: true, IsSigner: false},
		{PublicKey: m.toPool, IsWritable: true, IsSigner: false},
		{PublicKey: m.tokenAccount, IsWritable: true, IsSigner: false},
		{PublicKey: TokenProgramID, IsWritable: false, IsSigner: false},
	}

	return solana.NewInstruction(
		m.programID,
		accounts,
		data,
	), nil
}
