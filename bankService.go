package main

import "fmt"

type BankService interface {
	Transfer(fromAccount, toAccount BankAccount, amount float64, reference string) (BankAccount, BankAccount, error)
	GetStatement(account BankAccount) ([]Transaction, error)
}

type bankService struct {
	transactions map[AccountID]([]Transaction)
}

func NewBankService() BankService {
	return &bankService{
		transactions: make(map[AccountID]([]Transaction)),
	}
}

func (b *bankService) addTransaction(accountID AccountID, transaction Transaction) {
	list := b.transactions[accountID]
	if list == nil {
		list = []Transaction{transaction}
	} else {
		list = append(list, transaction)
	}
	b.transactions[accountID] = list
}

func (b *bankService) Transfer(fromAccount, toAccount BankAccount, amount float64, reference string) (BankAccount, BankAccount, error) {
	newFrom, err := fromAccount.Withdraw(amount)
	if err != nil {
		return nil, nil, fmt.Errorf("fromAccount: withdrawal failed: %w", err)
	}
	newTo, err := toAccount.Deposit(amount)
	if err != nil {
		return nil, nil, fmt.Errorf("toAccount: deposit failed: %w", err)
	}
	t := *NewTransaction(
		TransactionID(fmt.Sprintf("%s-%s", fromAccount.ID(), toAccount.ID())),
		fromAccount.ID(),
		toAccount.ID(),
		amount,
		reference,
	)
	b.addTransaction(fromAccount.ID(), t)
	b.addTransaction(toAccount.ID(), t)
	return newFrom, newTo, nil
}

func (b *bankService) GetStatement(account BankAccount) ([]Transaction, error) {
	if account == nil {
		return nil, fmt.Errorf("account is nil")
	}
	if b.transactions[account.ID()] == nil {
		// not an error because new accounts may not have been used yet
		// and therefore have no transactions
		return []Transaction{}, nil
	}
	return b.transactions[account.ID()], nil
}
