package main

import (
	"github.com/gagliardetto/solana-go"
)

// Transaction represents a parsed Solana transaction with Raydium-specific data
type Transaction struct {
	Signature solana.Signature
	Slot      uint64

	Create     []CreateInfo
	Trade      []TradeInfo
	TradeBuys  []int
	TradeSells []int

	Migrate   []Migration
	SwapBuys  []SwapBuy
	SwapSells []SwapSell
}

// CreateInfo represents token/pool creation information
type CreateInfo struct {
	TokenMint     solana.PublicKey
	TokenDecimals uint8
	TokenSymbol   string
	PoolAddress   solana.PublicKey
	Creator       solana.PublicKey
	Amount        uint64
	Timestamp     int64
}

// TradeInfo represents general trade information
type TradeInfo struct {
	InstructionIndex int
	TokenIn          solana.PublicKey
	TokenOut         solana.PublicKey
	AmountIn         uint64
	AmountOut        uint64
	Trader           solana.PublicKey
	Pool             solana.PublicKey
	TradeType        string // "buy", "sell", "swap"
}

// Migration represents a migration operation
type Migration struct {
	FromPool  solana.PublicKey
	ToPool    solana.PublicKey
	Token     solana.PublicKey
	Amount    uint64
	Owner     solana.PublicKey
	Timestamp int64
}

// SwapBuy represents a buy swap operation
type SwapBuy struct {
	TokenIn      solana.PublicKey
	TokenOut     solana.PublicKey
	AmountIn     uint64
	AmountOut    uint64
	MinAmountOut uint64
	Pool         solana.PublicKey
	Buyer        solana.PublicKey
	Slippage     float64
}

// SwapSell represents a sell swap operation
type SwapSell struct {
	TokenIn      solana.PublicKey
	TokenOut     solana.PublicKey
	AmountIn     uint64
	AmountOut    uint64
	MinAmountOut uint64
	Pool         solana.PublicKey
	Seller       solana.PublicKey
	Slippage     float64
}
