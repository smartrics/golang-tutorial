package bank_test

import (
	"testing"

	"github.com/smartrics/golang-tutorial/internal/bank"
)

func TestIntegration_Transfer_Success(t *testing.T) {
	service := bank.NewBankService()

	// Setup
	from := bank.NewBankAccount("ACC001", 1000)
	to := bank.NewBankAccount("ACC002", 500)

	// Act
	updatedFrom, updatedTo, err := service.Transfer(from, to, 200, "INV-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert balances
	if updatedFrom.Balance() != 800 {
		t.Errorf("expected updatedFrom balance 800, got %.2f", updatedFrom.Balance())
	}
	if updatedTo.Balance() != 700 {
		t.Errorf("expected updatedTo balance 700, got %.2f", updatedTo.Balance())
	}

	// Assert statement
	stmtFrom, err := service.GetStatement(updatedFrom)
	if err != nil {
		t.Fatalf("unexpected error from GetStatement(from): %v", err)
	}
	stmtTo, err := service.GetStatement(updatedTo)
	if err != nil {
		t.Fatalf("unexpected error from GetStatement(to): %v", err)
	}

	if len(stmtFrom) == 0 {
		t.Errorf("expected statement for from-account, got none")
	}
	if len(stmtTo) == 0 {
		t.Errorf("expected statement for to-account, got none")
	}
}
