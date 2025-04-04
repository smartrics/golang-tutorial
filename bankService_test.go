package main

import (
	"errors"
	"testing"
)

func TestBankService_Transfer(t *testing.T) {
	tests := []struct {
		name         string
		from         BankAccount
		to           BankAccount
		amount       float64
		expectErr    error
		expectedFrom float64
		expectedTo   float64
	}{
		{
			name:         "successful transfer",
			from:         NewBankAccount("123", 1000),
			to:           NewBankAccount("456", 500),
			amount:       200,
			expectErr:    nil,
			expectedFrom: 800,
			expectedTo:   700,
		},
		{
			name:      "withdrawal failure due to insufficient funds",
			from:      NewBankAccount("123", 100),
			to:        NewBankAccount("456", 500),
			amount:    200,
			expectErr: ErrFromAccountWithdrawal,
		},
		{
			name:      "invalid destination account",
			from:      NewBankAccount("123", 1000),
			to:        nil,
			amount:    200,
			expectErr: ErrInvalidAccount,
		},
		{
			name:      "self-transfer disallowed",
			from:      NewBankAccount("123", 1000),
			to:        nil, // set later as same as from
			amount:    200,
			expectErr: ErrSelfTransferDisallowed,
		},
		{
			name:         "transfer zero amount",
			from:         NewBankAccount("123", 1000),
			to:           NewBankAccount("456", 500),
			amount:       0,
			expectErr:    nil,
			expectedFrom: 1000,
			expectedTo:   500,
		},
		{
			name:      "transfer negative amount fails",
			from:      NewBankAccount("123", 1000),
			to:        NewBankAccount("456", 500),
			amount:    -200,
			expectErr: ErrFromAccountWithdrawal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bs := NewBankService()

			// Handle self-transfer case
			if tc.name == "self-transfer disallowed" {
				tc.to = tc.from
			}

			from, to, err := bs.Transfer(tc.from, tc.to, tc.amount, "ref")

			if tc.expectErr != nil {
				if err == nil {
					t.Errorf("Expected error '%v', got nil", tc.expectErr)
				} else if !errors.Is(err, tc.expectErr) {
					t.Errorf("Expected error '%v', got '%v'", tc.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if from.Balance() != tc.expectedFrom {
				t.Errorf("Expected from balance %.2f, got %.2f", tc.expectedFrom, from.Balance())
			}
			if to.Balance() != tc.expectedTo {
				t.Errorf("Expected to balance %.2f, got %.2f", tc.expectedTo, to.Balance())
			}
		})
	}
}

func TestBankService_GetStatement(t *testing.T) {
	bs := NewBankService()

	t.Run("returns error for nil account", func(t *testing.T) {
		_, err := bs.GetStatement(nil)
		if err == nil || !errors.Is(err, ErrInvalidAccount) {
			t.Errorf("Expected invalid account error, got %v", err)
		}
	})

	t.Run("returns empty for new account", func(t *testing.T) {
		acc := NewBankAccount("123", 1000)
		stmt, err := bs.GetStatement(acc)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(stmt) != 0 {
			t.Errorf("Expected empty statement, got %v", stmt)
		}
	})

	t.Run("returns statement after transfer", func(t *testing.T) {
		acc1 := NewBankAccount("123", 1000)
		acc2 := NewBankAccount("456", 500)
		_, _, _ = bs.Transfer(acc1, acc2, 200, "ref")

		stmt1, err := bs.GetStatement(acc1)
		if err != nil || len(stmt1) == 0 {
			t.Errorf("Expected statement for from-account, got %v, err: %v", stmt1, err)
		}

		stmt2, err := bs.GetStatement(acc2)
		if err != nil || len(stmt2) == 0 {
			t.Errorf("Expected statement for to-account, got %v, err: %v", stmt2, err)
		}
	})
}
