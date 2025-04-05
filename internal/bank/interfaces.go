package bank

type Identifier interface {
	ID() AccountID
}

type Balancer interface {
	Balance() float64
}

type Depositor interface {
	Deposit(amount float64) (BankAccount, error)
}

type Withdrawer interface {
	Withdraw(amount float64) (BankAccount, error)
}

type Stringer interface {
	String() string
}
