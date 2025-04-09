package ports

import "github.com/smartrics/golang-tutorial/internal/bank"

type BankServicePort interface {
	Transfer(from, to bank.BankAccount, amount float64, ref string) (bank.BankAccount, bank.BankAccount, error)
	GetStatement(bank.BankAccount) ([]bank.Transaction, error)
}
