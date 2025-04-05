package bank

type TransactionCounter struct {
	count int
}

func (tc *TransactionCounter) Increment() {
	tc.count++
}

func (tc *TransactionCounter) Count() int {
	return tc.count
}
