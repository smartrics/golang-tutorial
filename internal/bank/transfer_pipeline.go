package bank

import (
	"fmt"
	"log"
	"time"
)

var (
	ErrNilAccount = fmt.Errorf("%w: from and to accounts must not be nil", ErrInvalidAccount)
)

type TransferFunc func(from, to BankAccount, amount float64, reference string) (BankAccount, BankAccount, error)

func CoreTransfer(service BankService) TransferFunc {
	return func(from, to BankAccount, amount float64, reference string) (BankAccount, BankAccount, error) {
		return service.Transfer(from, to, amount, reference)
	}
}

func WithValidation(next TransferFunc) TransferFunc {
	return func(from, to BankAccount, amount float64, reference string) (BankAccount, BankAccount, error) {
		if from == nil || to == nil {
			return nil, nil, ErrNilAccount
		}
		if from.ID() == to.ID() {
			return nil, nil, ErrSelfTransferDisallowed
		}
		if amount < 0 {
			return nil, nil, ErrNegativeAmount
		}
		return next(from, to, amount, reference)
	}
}

func WithLogging(next TransferFunc) TransferFunc {
	return func(from, to BankAccount, amount float64, reference string) (BankAccount, BankAccount, error) {
		start := time.Now()
		log.Printf("[TRANSFER] initiating from: `%s`, to: `%s`, amount: `%.2f`, ref=`%s`", from.ID(), to.ID(), amount, reference)
		newFrom, newTo, err := next(from, to, amount, reference)
		if err != nil {
			log.Printf("[TRANSFER] FAILED after %s: %v", time.Since(start), err)
			return newFrom, newTo, err
		}

		log.Printf("[TRANSFER] SUCCESS after %s: newFromBalance=%.2f newToBalance=%.2f", time.Since(start), newFrom.Balance(), newTo.Balance())
		return newFrom, newTo, nil
	}
}
