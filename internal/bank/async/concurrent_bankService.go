package async

import (
	"sync"

	"github.com/smartrics/golang-tutorial/internal/bank"
)

// ConcurrentBankService wraps a bank.BankService to provide thread-safe access
// for use in concurrent environments like the async processor.
type ConcurrentBankService struct {
	inner bank.BankService
	mu    sync.Mutex
}

// NewConcurrentBankService returns a thread-safe wrapper around a new bankService.
func NewConcurrentBankService() *ConcurrentBankService {
	return &ConcurrentBankService{
		inner: bank.NewBankService(),
	}
}

// Transfer safely delegates to the inner BankService.Transfer with synchronized access.
func (c *ConcurrentBankService) Transfer(from, to bank.BankAccount, amount float64, ref string) (bank.BankAccount, bank.BankAccount, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.inner.Transfer(from, to, amount, ref)
}

// GetStatement safely delegates to the inner BankService.GetStatement with synchronized access.
func (c *ConcurrentBankService) GetStatement(account bank.BankAccount) ([]bank.Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.inner.GetStatement(account)
}
