package engine

import (
	"fmt"
	"sync"

	"github.com/smartrics/golang-tutorial/internal/bank"
)

var (
	ErrUnknownAccount = fmt.Errorf("unknown account")
)

// AccountRegistry provides thread-safe access to registered accounts.
type AccountRegistry interface {
	Get(id string) (bank.BankAccount, error)
	Register(bank.BankAccount)
	List() []bank.BankAccount
}

type accountRegistry struct {
	mu       sync.RWMutex
	accounts map[string]bank.BankAccount
}

// NewAccountRegistry creates a new empty account registry.
func NewAccountRegistry() AccountRegistry {
	return &accountRegistry{
		accounts: make(map[string]bank.BankAccount),
	}
}

func (r *accountRegistry) Register(acc bank.BankAccount) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.accounts[string(acc.ID())] = acc
}

func (r *accountRegistry) Get(id string) (bank.BankAccount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	acc, ok := r.accounts[id]
	if !ok {
		return nil, ErrUnknownAccount
	}
	return acc, nil
}

func (r *accountRegistry) List() []bank.BankAccount {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]bank.BankAccount, 0, len(r.accounts))
	for _, acc := range r.accounts {
		out = append(out, acc)
	}
	return out
}
