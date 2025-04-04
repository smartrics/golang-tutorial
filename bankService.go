package main

import "fmt"

var (
	ErrInvalidAccount         = fmt.Errorf("invalid account")
	ErrFromAccountWithdrawal  = fmt.Errorf("withdrawal failed from account")
	ErrToAccountDeposit       = fmt.Errorf("deposit failed to account")
	ErrSelfTransferDisallowed = fmt.Errorf("transfer from/to the same account is disallowed")
)

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
	if toAccount == nil || fromAccount == nil {
		return nil, nil, fmt.Errorf("%w: both accounts must be provided", ErrInvalidAccount)
	}
	if toAccount.ID() == fromAccount.ID() {
		return nil, nil, fmt.Errorf("%w: cannot transfer to/from the same account", ErrSelfTransferDisallowed)
	}

	newFrom, err := fromAccount.Withdraw(amount)
	if err != nil {
		return nil, nil, fmt.Errorf("%v: %w", err, ErrFromAccountWithdrawal)
	}
	newTo, err := toAccount.Deposit(amount)
	if err != nil {
		return nil, nil, fmt.Errorf("%v: %w", err, ErrToAccountDeposit)
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
		return nil, fmt.Errorf("%w: account must be provided", ErrInvalidAccount)
	}
	if b.transactions[account.ID()] == nil {
		// not an error because new accounts may not have been used yet
		// and therefore have no transactions
		return []Transaction{}, nil
	}
	return b.transactions[account.ID()], nil
}
