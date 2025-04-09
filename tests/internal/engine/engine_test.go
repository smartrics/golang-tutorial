package engine_test

import (
	"sync"
	"testing"
	"time"

	"github.com/smartrics/golang-tutorial/internal/bank"
	"github.com/smartrics/golang-tutorial/internal/bank/async"
	"github.com/smartrics/golang-tutorial/internal/engine"
)

func TestTransferEngine_SubmitTransfer_Succeeds(t *testing.T) {
	ctx := t.Context()

	// Create the engine and processor
	reg := engine.NewAccountRegistry()

	// Create and register accounts
	from := bank.NewBankAccount("A", 1000)
	to := bank.NewBankAccount("B", 0)

	reg.Register(from)
	reg.Register(to)

	// Hook for completion
	done := make(chan error, 1)

	// Wrap processor with audit handler
	var auditMu sync.Mutex
	var audit []bank.AuditEntry

	auditFn := func(e bank.AuditEntry) {
		auditMu.Lock()
		defer auditMu.Unlock()
		audit = append(audit, e)
	}

	svc := async.NewConcurrentBankService()
	transfer := bank.WithAudit(auditFn)(
		bank.WithValidation(
			bank.CoreTransfer(svc),
		),
	)
	proc := async.NewProcessor(transfer, 1)
	_ = proc.Start(ctx)

	// Create engine
	eng := engine.NewEngine(reg, proc, svc)
	eng.OnComplete(func(job async.TransferJob, err error) {
		done <- err
	})

	// Submit
	err := eng.SubmitTransfer("A", "B", 100, "ref-1")
	if err != nil {
		t.Fatalf("unexpected submit error: %v", err)
	}

	// Wait for job completion
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("transfer failed: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for job")
	}

	// Verify statement for 'A' and 'B'
	stA, err := eng.GetStatement("A")
	if err != nil {
		t.Fatalf("get statement A failed: %v", err)
	}
	stB, err := eng.GetStatement("B")
	if err != nil {
		t.Fatalf("get statement B failed: %v", err)
	}
	if len(stA) != 1 || len(stB) != 1 {
		t.Errorf("expected 1 transaction each, got %d and %d", len(stA), len(stB))
	}
}
