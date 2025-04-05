package bank

type AccountWithCounter struct {
	BankAccount
	*TransactionCounter
}

func NewAccountWithCounter(inner BankAccount) *AccountWithCounter {
	return newAccountWithCounterSet(inner, 0)
}

func newAccountWithCounterSet(inner BankAccount, txCount int) *AccountWithCounter {
	return &AccountWithCounter{
		BankAccount: inner,
		TransactionCounter: &TransactionCounter{
			count: txCount,
		},
	}
}

func (a *AccountWithCounter) Deposit(amount float64) (BankAccount, error) {
	newAcc, err := a.BankAccount.Deposit(amount)
	if err == nil {
		a.TransactionCounter.Increment()
	}
	return newAccountWithCounterSet(newAcc, a.count), err
}

func (a *AccountWithCounter) Withdraw(amount float64) (BankAccount, error) {
	newAcc, err := a.BankAccount.Withdraw(amount)
	if err == nil {
		a.TransactionCounter.Increment()
	}
	return newAccountWithCounterSet(newAcc, a.count), err
}
