package main

import (
	"errors"
	"testing"
)

func TestBankService_Transfer_WithMoqWithdrawSuccess(t *testing.T) {
	bs := NewBankService()

	idFuncFrom := func() AccountID { return "A" }
	balFuncFrom := func() float64 { return 1000 }
	from := &BankAccountMock{
		IDFunc:      idFuncFrom,
		BalanceFunc: balFuncFrom,
		WithdrawFunc: func(amount float64) (BankAccount, error) {
			return &BankAccountMock{
				IDFunc:      idFuncFrom,
				BalanceFunc: balFuncFrom,
			}, nil
		},
	}

	idFuncTo := func() AccountID { return "B" }
	balFuncTo := func() float64 { return 600 }
	to := &BankAccountMock{
		IDFunc:      idFuncTo,
		BalanceFunc: balFuncTo,
		DepositFunc: func(amount float64) (BankAccount, error) {
			return &BankAccountMock{
				IDFunc:      idFuncTo,
				BalanceFunc: balFuncTo,
			}, nil
		},
	}

	newFrom, newTo, err := bs.Transfer(from, to, 200, "ref")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if newFrom.Balance() != balFuncFrom() {
		t.Errorf("Expected from balance %.2f, got %.2f", balFuncFrom(), newFrom.Balance())
	}
	if newTo.Balance() != balFuncTo() {
		t.Errorf("Expected to balance %.2f, got %.2f", balFuncTo(), newTo.Balance())
	}
}

func TestBankService_Transfer_WithMoqWithdrawFailure(t *testing.T) {
	bs := NewBankService()

	from := &BankAccountMock{
		IDFunc:      func() AccountID { return "A" },
		BalanceFunc: func() float64 { return 100 },
		WithdrawFunc: func(amount float64) (BankAccount, error) {
			return nil, ErrFromAccountWithdrawal
		},
	}

	idFunc := func() AccountID { return "B" }
	balFunc := func() float64 { return 500 }
	to := &BankAccountMock{
		IDFunc:      idFunc,
		BalanceFunc: balFunc,
		DepositFunc: func(amount float64) (BankAccount, error) {
			return &BankAccountMock{
				IDFunc:      idFunc,
				BalanceFunc: balFunc,
			}, nil
		},
	}

	_, _, err := bs.Transfer(from, to, 200, "ref")

	if !errors.Is(err, ErrFromAccountWithdrawal) {
		t.Errorf("Expected ErrFromAccountWithdrawal, got %v", err)
	}
}
