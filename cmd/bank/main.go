package main

import "github.com/smartrics/golang-tutorial/internal/bank"

func main() {
	// Create a new transaction
	tx := bank.NewTransaction("123", "Alice", "Bob", 100.0, "test-ref")

	// Print transaction details
	println("Transaction ID:", tx.ID())
	println("From:", tx.From())
	println("To:", tx.To())
	println("Amount:", tx.Amount())
	println("Reference:", tx.Reference())
}
