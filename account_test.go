package main

import (
	"errors"
	"testing"
)

func TestNewAccountRequiresIDAndAmount(t *testing.T) {
	acc := NewBankAccount("ABC-1", 1000.0)
	if acc.ID() != "ABC-1" {
		t.Errorf("Expected account ID 'ABC-1', got '%s'", acc.ID())
	}
	if acc.Balance() != 1000.0 {
		t.Errorf("Expected account balance 1000.0, got '%.2f'", acc.Balance())
	}
}

func TestDepositPositiveAmountSucceeds(t *testing.T) {
	acc := NewBankAccount("ABC-1", 1000.0)
	acc, err := acc.Deposit(500.0)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	if acc.Balance() != 1500.0 {
		t.Errorf("Expected account balance 1500.0, got '%.2f'", acc.Balance())
	}
}

func TestWithdrawPositiveAmountSucceeds(t *testing.T) {
	acc := NewBankAccount("ABC-1", 1000.0)
	acc, err := acc.Withdraw(500.0)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	if acc.Balance() != 500.0 {
		t.Errorf("Expected account balance 500.0, got '%.2f'", acc.Balance())
	}
}

func TestDepositNegativeAmountFails(t *testing.T) {
	acc := NewBankAccount("ABC-1", 1000.0)
	acc, err := acc.Deposit(-500.0)
	if err == nil {
		if acc.Balance() != 1000.0 {
			t.Errorf("Expected account balance unchanged in case of error, had 1000.0 but got '%.2f'", acc.Balance())
		}
		t.Errorf("Expected an error")
	}
	if !errors.Is(err, ErrDepositNegativeAmount) {
		t.Errorf("Expected error 'cannot deposit negative amount', got '%v'", err)
	}
}

func TestWithdrawNegativeAmountFails(t *testing.T) {
	acc := NewBankAccount("ABC-1", 1000.0)
	acc, err := acc.Withdraw(-500.0)
	if err == nil {
		if acc.Balance() != 1000.0 {
			t.Errorf("Expected account balance unchanged in case of error, had 1000.0 but got '%.2f'", acc.Balance())
		}
		t.Errorf("Expected an error")
	}
	if !errors.Is(err, ErrWithdrawNegativeAmount) {
		t.Errorf("Expected error 'cannot withdraw negative amount', got '%v'", err)
	}
}

func TestWithdrawInsufficientFundsFails(t *testing.T) {
	acc := NewBankAccount("ABC-1", 1000.0)
	acc, err := acc.Withdraw(1500.0)
	if err == nil {
		if acc.Balance() != 1000.0 {
			t.Errorf("Expected account balance unchanged in case of error, had 1000.0 but got '%.2f'", acc.Balance())
		}
		t.Errorf("Expected an error")
	}
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("Expected error 'insufficient funds', got '%v'", err)
	}
}

func TestNewCheckingAccount(t *testing.T) {
	ca, err := NewCheckingAccount("ABC-1", 1000.0, 100.0)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	if ca.OverdraftLimit() != 100.0 {
		t.Errorf("Expected overdraft limit 100.0, got '%.2f'", ca.OverdraftLimit())
	}
}

func TestNewCheckingAccountNegativeAmountsWithinOverdraftSucceeds(t *testing.T) {
	ca, err := NewCheckingAccount("ABC-1", -50.0, 100.0)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	if ca.Balance() != -50.0 {
		t.Errorf("Expected balance -50.0, got '%.2f'", ca.Balance())
	}
}

func TestNewCheckingAccountWithBalanceLowerThanOverdraftFails(t *testing.T) {
	ca, err := NewCheckingAccount("ABC-1", -200.0, 100.0)
	if err == nil {
		t.Errorf("Expected error")
	}
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("Expected error 'insufficient funds', got '%v'", err)
	}
	if ca != nil {
		t.Errorf("Expected nil account, got '%v'", ca)
	}
}

func TestNewCheckingAccountWithdrawWithinOverdraftSucceeds(t *testing.T) {
	ca, _ := NewCheckingAccount("ABC-1", 1000.0, 100.0)
	ca2, _ := ca.Withdraw(1050.0)
	if ca2.Balance() != -50.0 {
		t.Errorf("Expected account balance -50.0, got '%.2f'", ca2.Balance())
	}
}

func TestNewCheckingAccountInvalidOverdraftLimit(t *testing.T) {
	ca, err := NewCheckingAccount("ABC-1", 1000.0, -100.0)
	if err == nil {
		t.Errorf("Expected error")
	}
	if ca != nil {
		t.Errorf("Expected nil account, got '%v'", ca)
	}
	if !errors.Is(err, ErrInvalidOverdraftLimit) {
		t.Errorf("Expected error 'invalid overdraft limit', got '%v'", err)
	}
}
