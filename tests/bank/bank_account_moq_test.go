package bank_test

import (
	"errors"
	"testing"

	"github.com/smartrics/golang-tutorial/internal/bank"
	bank_mocks "github.com/smartrics/golang-tutorial/mocks/bank"
)

func TestBankService_Transfer_WithMoqWithdrawSuccess(t *testing.T) {
	bs := bank.NewBankService()

	idFuncFrom := func() bank.AccountID { return "A" }
	balFuncFrom := func() float64 { return 1000 }
	from := &bank_mocks.BankAccountMock{
		IDFunc:      idFuncFrom,
		BalanceFunc: balFuncFrom,
		WithdrawFunc: func(amount float64) (bank.BankAccount, error) {
			return &bank_mocks.BankAccountMock{
				IDFunc:      idFuncFrom,
				BalanceFunc: balFuncFrom,
			}, nil
		},
	}

	idFuncTo := func() bank.AccountID { return "B" }
	balFuncTo := func() float64 { return 600 }
	to := &bank_mocks.BankAccountMock{
		IDFunc:      idFuncTo,
		BalanceFunc: balFuncTo,
		DepositFunc: func(amount float64) (bank.BankAccount, error) {
			return &bank_mocks.BankAccountMock{
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
	bs := bank.NewBankService()

	from := &bank_mocks.BankAccountMock{
		IDFunc:      func() bank.AccountID { return "A" },
		BalanceFunc: func() float64 { return 100 },
		WithdrawFunc: func(amount float64) (bank.BankAccount, error) {
			return nil, bank.ErrFromAccountWithdrawal
		},
	}

	idFunc := func() bank.AccountID { return "B" }
	balFunc := func() float64 { return 500 }
	to := &bank_mocks.BankAccountMock{
		IDFunc:      idFunc,
		BalanceFunc: balFunc,
		DepositFunc: func(amount float64) (bank.BankAccount, error) {
			return &bank_mocks.BankAccountMock{
				IDFunc:      idFunc,
				BalanceFunc: balFunc,
			}, nil
		},
	}

	_, _, err := bs.Transfer(from, to, 200, "ref")

	if !errors.Is(err, bank.ErrFromAccountWithdrawal) {
		t.Errorf("Expected bank.ErrFromAccountWithdrawal, got %v", err)
	}
}
