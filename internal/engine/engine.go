package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartrics/golang-tutorial/internal/bank"
	"github.com/smartrics/golang-tutorial/internal/bank/async"
)

type contextKey string

type AuditEntry struct {
	From      string
	To        string
	Amount    float64
	Reference string
	Timestamp time.Time
	Success   bool
	Error     string
}

var (
	auditMu  sync.Mutex
	auditLog []AuditEntry
)

const (
	contextKeyTimestamp = contextKey("ts")
	contextKeyJobID     = contextKey("jobID")
)

type BankServicePort interface {
	Transfer(from, to bank.BankAccount, amount float64, ref string) (bank.BankAccount, bank.BankAccount, error)
	GetStatement(bank.BankAccount) ([]bank.Transaction, error)
}

type TransferEngine interface {
	SubmitTransfer(fromID, toID string, amount float64, reference string) error
	GetStatement(accountID string) ([]bank.Transaction, error)
	OnComplete(func(job async.TransferJob, err error))
	RegisterAccount(bank.BankAccount)
	AuditLog(out *[]AuditEntry)
}

type transferEngine struct {
	registry  AccountRegistry
	processor async.Processor
	bankSvc   BankServicePort
	onDone    func(job async.TransferJob, err error)
}

func (e *transferEngine) AuditLog(out *[]AuditEntry) {
	if out == nil {
		panic("AuditLog: nil output slice")
	}

	auditMu.Lock()
	defer auditMu.Unlock()

	// Allocate a new slice or extend the given one
	copied := make([]AuditEntry, len(auditLog))
	copy(copied, auditLog)

	*out = copied
}

// New creates a new TransferEngine instance with the given registry and processor.
func NewEngine(reg AccountRegistry, proc async.Processor, svc BankServicePort) TransferEngine {
	return &transferEngine{
		registry:  reg,
		bankSvc:   svc,
		processor: proc,
	}
}

func (e *transferEngine) RegisterAccount(acc bank.BankAccount) {
	e.registry.Register(acc)
}

func (e *transferEngine) SubmitTransfer(fromID, toID string, amount float64, reference string) error {
	from, err := e.registry.Get(fromID)
	if err != nil {
		return fmt.Errorf("from: %w", err)
	}
	to, err := e.registry.Get(toID)
	if err != nil {
		return fmt.Errorf("to: %w", err)
	}
	if from.ID() == to.ID() {
		return fmt.Errorf("cannot transfer to self")
	}
	if amount <= 0 {
		return fmt.Errorf("invalid amount: %f", amount)
	}

	done := make(chan error, 1)

	job := async.TransferJob{
		From:   from,
		To:     to,
		Amount: amount,
		Ctx:    context.WithValue(context.Background(), contextKeyTimestamp, time.Now()),
		Done:   done,
	}

	err = e.processor.SubmitWithRetry(job, 3, 100*time.Millisecond)
	if err != nil {
		return fmt.Errorf("submit: %w", err)
	}

	go func() {
		result := <-done
		if e.onDone != nil {
			e.onDone(job, result)
		}
		auditMu.Lock()
		auditLog = append(auditLog, AuditEntry{
			From:      fromID,
			To:        toID,
			Amount:    amount,
			Reference: reference,
			Timestamp: time.Now(),
			Success:   result == nil,
			Error:     errString(result),
		})
		auditMu.Unlock()
	}()

	return nil
}

func (e *transferEngine) GetStatement(accountID string) ([]bank.Transaction, error) {
	acc, err := e.registry.Get(accountID)
	if err != nil {
		return nil, fmt.Errorf("statement: %w", err)
	}
	return e.bankSvc.GetStatement(acc)
}

func (e *transferEngine) OnComplete(cb func(job async.TransferJob, err error)) {
	e.onDone = cb
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
