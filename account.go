package main

import (
	"fmt"
)

var (
	ErrWithdrawNegativeAmount = fmt.Errorf("cannot withdraw negative amount")
	ErrDepositNegativeAmount  = fmt.Errorf("cannot deposit negative amount")
	ErrInsufficientFunds      = fmt.Errorf("insufficient funds")
	ErrInvalidOverdraftLimit  = fmt.Errorf("invalid overdraft limit")
	ErrInvalidInterestRate    = fmt.Errorf("invalid interest rate")
)

type AccountID string

type BankAccount interface {
	Withdraw(amount float64) (BankAccount, error)
	Deposit(amount float64) (BankAccount, error)
	Balance() float64
	String() string
	ID() AccountID
}

type SavingsAccount interface {
	BankAccount
	InterestRate() float64
}
type CheckingAccount interface {
	BankAccount
	OverdraftLimit() float64
}

var _ SavingsAccount = &savingsAccount{}
var _ CheckingAccount = &checkingAccount{}

type bankAccount struct {
	id      AccountID
	balance float64
}

type savingsAccount struct {
	bankAccount
	interestRate float64
}

type checkingAccount struct {
	bankAccount
	overdraftLimit float64
}

func NewBankAccount(id AccountID, balance float64) BankAccount {
	return bankAccount{id: id, balance: balance}
}

func NewCheckingAccount(id AccountID, balance float64, overdraftLimit float64) (CheckingAccount, error) {
	if overdraftLimit < 0 {
		return nil, ErrInvalidOverdraftLimit
	}
	if balance < -overdraftLimit {
		return nil, ErrInsufficientFunds
	}
	ba := bankAccount{id: id, balance: balance}
	return checkingAccount{bankAccount: ba, overdraftLimit: overdraftLimit}, nil
}

func NewSavingAccount(id AccountID, balance float64, interestRate float64) (SavingsAccount, error) {
	if interestRate < 0 {
		return nil, ErrInvalidInterestRate
	}
	ba := bankAccount{id: id, balance: balance}
	return savingsAccount{bankAccount: ba, interestRate: 0}, nil
}

func (a bankAccount) Withdraw(amount float64) (BankAccount, error) {
	if amount < 0 {
		return a, ErrWithdrawNegativeAmount
	}
	if amount > a.balance {
		return a, ErrInsufficientFunds
	}
	return NewBankAccount(a.id, a.balance-amount), nil
}

func (a bankAccount) Deposit(amount float64) (BankAccount, error) {
	if amount < 0 {
		return a, ErrDepositNegativeAmount
	}
	return NewBankAccount(a.id, a.balance+amount), nil
}

func (a bankAccount) Balance() float64 {
	return a.balance
}

func (a bankAccount) String() string {
	return fmt.Sprintf("Account ID: %s, Balance: %.2f", a.id, a.balance)
}

func (a bankAccount) ID() AccountID {
	return a.id
}

// OverdraftLimit implements CheckingAccount.
func (c checkingAccount) OverdraftLimit() float64 {
	return c.overdraftLimit
}

// Withdraw implements CheckingAccount.
// Subtle: this method shadows the method (bankAccount).Withdraw of checkingAccount.bankAccount.
func (c checkingAccount) Withdraw(amount float64) (BankAccount, error) {
	return NewCheckingAccount(c.id, c.balance-amount, c.overdraftLimit)
}

// InterestRate implements SavingsAccount.
func (s savingsAccount) InterestRate() float64 {
	return s.interestRate
}
