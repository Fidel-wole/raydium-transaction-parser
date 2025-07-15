package main

import (
	"fmt"
	"log"
	"testing"
)

func TestSpecificTransaction(t *testing.T) {

	signature := "2N9VyxzFmHibuWy5HmJH52R6Hy6NZPw5iCdFc9X1JT4JBPCa4VZmxv3RhSvP9UfDdCdgDYvoeaN62v29toJNAWtD"

	fmt.Printf("Testing specific transaction: %s\n", signature)

	log.Printf("Expected account mapping for Launchpad create:")
	log.Printf("  Account[0] = Creator/Signer")
	log.Printf("  Account[1] = Token Mint (new token)")
	log.Printf("  Account[2] = Pool Address")
	log.Printf("  Account[3+] = Various system/program accounts")

	log.Printf("Expected account mapping for Launchpad buy:")
	log.Printf("  Account[0] = Buyer/Signer")
	log.Printf("  Account[1] = Token Mint (token being bought)")
	log.Printf("  Account[2] = Pool Address")
	log.Printf("  Account[3+] = Token accounts and system programs")
}

func TestCurrentParsingLogic(t *testing.T) {
	fmt.Println("Testing current parsing logic...")

	fmt.Println("Current issue:")
	fmt.Println("- Parser is using instruction.Accounts[0] as token mint")
	fmt.Println("- But instruction.Accounts[0] is actually the creator/signer")
	fmt.Println("- Token mint should be instruction.Accounts[1]")
	fmt.Println("- Pool should be instruction.Accounts[2]")

	fmt.Println("\nCorrect mapping should be:")
	fmt.Println("- Creator = message.AccountKeys[0] (signer)")
	fmt.Println("- Token Mint = message.AccountKeys[instruction.Accounts[1]]")
	fmt.Println("- Pool = message.AccountKeys[instruction.Accounts[2]]")
}
