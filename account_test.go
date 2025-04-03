package main

import (
	"errors"
	"testing"
)

func TestNewAccountRequiresIDAndAmount(t *testing.T) {
	acc := NewAccount("ABC-1", 1000.0)
	if acc.id != "ABC-1" {
		t.Errorf("Expected account ID 'ABC-1', got '%s'", acc.id)
	}
	if acc.balance != 1000.0 {
		t.Errorf("Expected account balance 1000.0, got '%.2f'", acc.balance)
	}
}

func TestDepositPositiveAmountSucceeds(t *testing.T) {
	acc := NewAccount("ABC-1", 1000.0)
	acc, err := acc.DepositIM(500.0)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	if acc.balance != 1500.0 {
		t.Errorf("Expected account balance 1500.0, got '%.2f'", acc.balance)
	}
}

func TestWithdrawPositiveAmountSucceeds(t *testing.T) {
	acc := NewAccount("ABC-1", 1000.0)
	acc, err := acc.WithdrawIM(500.0)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	if acc.balance != 500.0 {
		t.Errorf("Expected account balance 500.0, got '%.2f'", acc.balance)
	}
}

func TestDepositNegativeAmountFails(t *testing.T) {
	acc := NewAccount("ABC-1", 1000.0)
	acc, err := acc.DepositIM(-500.0)
	if err == nil {
		if acc.balance != 1000.0 {
			t.Errorf("Expected account balance unchanged in case of error, had 1000.0 but got '%.2f'", acc.balance)
		}
		t.Errorf("Expected an error")
	}
	if !errors.Is(err, ErrDepositNegativeAmount) {
		t.Errorf("Expected error 'cannot deposit negative amount', got '%v'", err)
	}
}

func TestWithdrawNegativeAmountFails(t *testing.T) {
	acc := NewAccount("ABC-1", 1000.0)
	acc, err := acc.WithdrawIM(-500.0)
	if err == nil {
		if acc.balance != 1000.0 {
			t.Errorf("Expected account balance unchanged in case of error, had 1000.0 but got '%.2f'", acc.balance)
		}
		t.Errorf("Expected an error")
	}
	if !errors.Is(err, ErrWithdrawNegativeAmount) {
		t.Errorf("Expected error 'cannot withdraw negative amount', got '%v'", err)
	}
}

func TestWithdrawInsufficientFundsFails(t *testing.T) {
	acc := NewAccount("ABC-1", 1000.0)
	acc, err := acc.WithdrawIM(1500.0)
	if err == nil {
		if acc.balance != 1000.0 {
			t.Errorf("Expected account balance unchanged in case of error, had 1000.0 but got '%.2f'", acc.balance)
		}
		t.Errorf("Expected an error")
	}
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("Expected error 'insufficient funds', got '%v'", err)
	}
}
