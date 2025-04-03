package main

import (
	"fmt"
)

var (
	ErrWithdrawNegativeAmount = fmt.Errorf("cannot withdraw negative amount")
	ErrDepositNegativeAmount  = fmt.Errorf("cannot deposit negative amount")
	ErrInsufficientFunds      = fmt.Errorf("insufficient funds")
)

type Account struct {
	id      string
	balance float64
}

func NewAccount(id string, balance float64) Account {
	return Account{id: id, balance: balance}
}

func (a Account) WithdrawIM(amount float64) (Account, error) {
	if amount < 0 {
		return a, ErrWithdrawNegativeAmount
	}
	if amount > a.balance {
		return a, ErrInsufficientFunds
	}
	return NewAccount(a.id, a.balance-amount), nil
}

func (a Account) DepositIM(amount float64) (Account, error) {
	if amount < 0 {
		return a, ErrDepositNegativeAmount
	}
	return NewAccount(a.id, a.balance+amount), nil
}

func (a Account) Balance() float64 {
	return a.balance
}

func (a Account) String() string {
	return fmt.Sprintf("Account ID: %s, Balance: %.2f", a.id, a.balance)
}

func main() {
	account := NewAccount("ABC-1", 1000.0)
	fmt.Println(account)
	account, _ = account.DepositIM(500.0)
	fmt.Println(account)
	account, _ = account.WithdrawIM(700.0)
	fmt.Println(account)
	fmt.Println("---")

	account2 := NewAccount("ABC-2", 2000.0)
	account3 := NewAccount("ABC-3", 2000.0)

	accounts := append([]Account{account, account2}, account3)
	for _, acc := range accounts {
		accountUpdated, _ := acc.DepositIM(100.0)
		fmt.Println(accountUpdated)
	}

	fmt.Println("---")
	accounts = append([]Account{account, account2}, account3)
	for _, acc := range accounts {
		fmt.Println(acc)
	}
}
