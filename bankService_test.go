package main

import (
	"testing"
)

func TestBankServiceTransferUpdatesAccounts(t *testing.T) {
	bs := NewBankService()
	account1 := NewBankAccount("123", 1000)
	account2 := NewBankAccount("456", 500)
	account1, account2, _ = bs.Transfer(account1, account2, 200, "ref")
	if account1.Balance() != 800 {
		t.Errorf("Expected account1 balance to be 800, got %f", account1.Balance())
	}
	if account2.Balance() != 700 {
		t.Errorf("Expected account2 balance to be 700, got %f", account2.Balance())
	}
}

func TestBankServiceTransferWithdrawalFailure(t *testing.T) {
	bs := NewBankService()
	account1 := NewBankAccount("123", 100)
	account2 := NewBankAccount("456", 500)
	_, _, err := bs.Transfer(account1, account2, 200, "ref")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if err.Error() != "fromAccount: withdrawal failed: insufficient funds" {
		t.Errorf("Expected insufficient funds error, got %v", err)
	}
}

func TestTransferGeneratesStatementForFromAndToAccounts(t *testing.T) {
	bs := NewBankService()
	account1 := NewBankAccount("123", 1000)
	account2 := NewBankAccount("456", 500)
	_, _, _ = bs.Transfer(account1, account2, 200, "ref")

	statement1, err := bs.GetStatement(account1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(statement1) == 0 {
		t.Errorf("Expected a transaction for account1, got %v", statement1)
	}

	statement2, err := bs.GetStatement(account2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(statement2) == 0 {
		t.Errorf("Expected a transaction for account2, got %v", statement2)
	}
}

func TestGetStatementReturnsEmptyForNewAccount(t *testing.T) {
	bs := NewBankService()
	account1 := NewBankAccount("123", 1000)
	statement, err := bs.GetStatement(account1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(statement) != 0 {
		t.Errorf("Expected empty statement for new account, got %v", statement)
	}
}
