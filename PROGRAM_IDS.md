# Raydium Program IDs Used in Testing

Based on my analysis of the Raydium Parser codebase, here are all the program IDs currently being used:

## **Primary Raydium Program IDs**

### 1. **Raydium V4 AMM Program** 
- **Program ID**: `675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8`
- **Used for**: Swap instructions and Migrate instructions
- **Location**: `parser.go` line 15, `instructions.go` lines 36, 590

### 2. **Raydium V5 AMM Program**
- **Program ID**: `5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h`
- **Used for**: Advanced swap operations
- **Location**: `parser.go` line 16

### 3. **Raydium Launchpad V1 Program**
- **Program ID**: `6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P`
- **Used for**: Buy instructions, Sell instructions, and CreateToken instructions
- **Location**: `parser.go` line 20, `instructions.go` lines 215, 332, 448

### 4. **Raydium CP-Swap Program**
- **Program ID**: `CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C`
- **Used for**: Concentrated pool swaps
- **Location**: `parser.go` line 21

## **Additional Raydium Program IDs**

### 5. **Raydium Staking Program**
- **Program ID**: `EhhTKczWMGQt46ynNeRX1WfeagwwJd7ufHvCDjRxjo5Q`
- **Used for**: Staking operations
- **Location**: `parser.go` line 17

### 6. **Raydium Liquidity Program**
- **Program ID**: `27haf8L6oxUeXrHrgEgsexjSY5hbVUWEmvv9Nyxg8vQv`
- **Used for**: Liquidity pool operations
- **Location**: `parser.go` line 18

### 7. **Unknown Raydium Program IDs** (found in real transactions)
- **Program ID 1**: `FoaFt2Dtz58RA6DPjbRb9t9z8sLJRChiGFTv21EfaseZ`
- **Program ID 2**: `LanMV9sAd7wArD4vJFi2qDdfnVhFxYSUg6eADduJ3uj`
- **Used for**: Generic Raydium instruction parsing
- **Location**: `parser.go` lines 23-24

## **Supporting Solana Program IDs**

### 8. **Token Program**
- **Program ID**: `TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA`
- **Used for**: Token operations in instruction builders
- **Location**: `parser.go` line 26, `instructions.go` lines 187, 303

### 9. **Token-2022 Program**
- **Program ID**: `TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb`
- **Used for**: Token-2022 operations
- **Location**: `parser.go` line 27

### 10. **System Program**
- **Program ID**: `11111111111111111111111111111111`
- **Used for**: System operations in instruction builders
- **Location**: `parser.go` line 28, `instructions.go` line 304

### 11. **Associated Token Program**
- **Program ID**: `ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL`
- **Used for**: Associated token account operations
- **Location**: `parser.go` line 29

## **Instruction Builder Program ID Mapping**

| Instruction Type | Default Program ID | Can be Changed? |
|------------------|-------------------|----------------|
| SwapInstruction | `675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8` (V4) | ✅ Yes |
| BuyInstruction | `6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P` (Launchpad) | ✅ Yes |
| SellInstruction | `6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P` (Launchpad) | ✅ Yes |
| CreateTokenInstruction | `6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P` (Launchpad) | ✅ Yes |
| MigrateInstruction | `675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8` (V4) | ✅ Yes |

## **Test Results**

When running the instruction builder tests (`go run . test`), you'll see output like:
```
1. Testing Swap Instruction Builder:
   - Program ID: 675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8

2. Testing Buy Instruction Builder:
   - Program ID: 6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P

3. Testing Sell Instruction Builder:
   - Program ID: 6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P

4. Testing Create Token Instruction Builder:
   - Program ID: 6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P

5. Testing Migrate Instruction Builder:
   - Program ID: 675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8
```

## **Summary**

The codebase primarily uses **two main program IDs** for testing:
1. **Raydium V4 AMM**: `675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8` (for swaps and migrations)
2. **Raydium Launchpad V1**: `6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P` (for buy/sell/create operations)

All instruction builders allow you to change the program ID using the `SetProgramID()` method if needed for different Raydium program versions or testing scenarios.
