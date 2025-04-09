package async_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/smartrics/golang-tutorial/internal/bank"
	"github.com/smartrics/golang-tutorial/internal/bank/async"
)

func TestProcessorExecutesTransferJob(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup fake service with audit to capture execution
	service := bank.NewBankService()
	var audit []bank.AuditEntry
	auditFn := func(e bank.AuditEntry) {
		audit = append(audit, e)
	}

	// Build pipeline
	pipeline := bank.WithAudit(auditFn)(
		bank.WithValidation(
			bank.CoreTransfer(service),
		),
	)

	// Create processor with N=1 worker
	proc := async.NewProcessor(pipeline, 1)
	if err := proc.Start(ctx); err != nil {
		t.Fatalf("failed to start processor: %v", err)
	}

	// Create accounts and job
	from := bank.NewBankAccount("FROM", 1000)
	to := bank.NewBankAccount("TO", 500)

	done := make(chan error, 1)
	job := async.TransferJob{
		From:      from,
		To:        to,
		Amount:    250,
		Reference: "REF-123",
		Done:      done,
	}

	// Submit job
	if err := proc.Submit(job); err != nil {
		t.Fatalf("failed to submit job: %v", err)
	}

	// Wait for result
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("transfer failed: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("transfer did not complete in time")
	}

	// Validate audit captured
	if len(audit) != 1 {
		t.Errorf("expected 1 audit entry, got %d", len(audit))
	}
}

func TestProcessor_ConcurrentJobs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := bank.NewBankService()
	var audit []bank.AuditEntry
	auditFn := func(e bank.AuditEntry) {
		audit = append(audit, e)
	}

	transfer := bank.WithAudit(auditFn)(
		bank.WithValidation(
			bank.CoreTransfer(service),
		),
	)

	proc := async.NewProcessor(transfer, 4) // multiple workers
	if err := proc.Start(ctx); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	const numJobs = 10
	var wg sync.WaitGroup
	wg.Add(numJobs)

	for i := range numJobs {
		from := bank.NewBankAccount(bank.AccountID(fmt.Sprintf("FROM-%d", i)), 1000)
		to := bank.NewBankAccount(bank.AccountID(fmt.Sprintf("TO-%d", i)), 0)

		job := async.TransferJob{
			From:      from,
			To:        to,
			Amount:    100,
			Reference: fmt.Sprintf("JOB-%d", i),
			Done:      make(chan error, 1),
		}

		go func(j async.TransferJob) {
			defer wg.Done()
			if err := proc.Submit(j); err != nil {
				t.Errorf("submit error: %v", err)
			}
			select {
			case err := <-j.Done:
				if err != nil {
					t.Errorf("transfer failed: %v", err)
				}
			case <-time.After(1 * time.Second):
				t.Error("timeout waiting for job completion")
			}
		}(job)
	}

	wg.Wait()

	if len(audit) != numJobs {
		t.Errorf("expected %d audit entries, got %d", numJobs, len(audit))
	}
}

func TestProcessor_StopShutsDownWorkers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Track executed jobs
	var mu sync.Mutex
	var executed []string

	transfer := func(from, to bank.BankAccount, amount float64, ref string) (bank.BankAccount, bank.BankAccount, error) {
		time.Sleep(50 * time.Millisecond) // simulate work
		mu.Lock()
		executed = append(executed, ref)
		mu.Unlock()
		return to, from, nil
	}

	proc := async.NewProcessor(transfer, 2)
	if err := proc.Start(ctx); err != nil {
		t.Fatalf("failed to start processor: %v", err)
	}

	// Submit 2 jobs
	for i := range 2 {
		from := bank.NewBankAccount(bank.AccountID(fmt.Sprintf("F%d", i)), 100)
		to := bank.NewBankAccount(bank.AccountID(fmt.Sprintf("T%d", i)), 100)
		job := async.TransferJob{
			From:      from,
			To:        to,
			Amount:    50,
			Reference: fmt.Sprintf("REF-%d", i),
			Done:      make(chan error, 1),
		}
		if err := proc.Submit(job); err != nil {
			t.Fatalf("failed to submit: %v", err)
		}
	}
	time.Sleep(100 * time.Millisecond)
	// Graceful stop
	if err := proc.Stop(); err != nil {
		t.Errorf("stop error: %v", err)
	}

	// Try submitting after stop (should fail)
	from := bank.NewBankAccount("F99", 100)
	to := bank.NewBankAccount("T99", 100)
	job := async.TransferJob{
		From:      from,
		To:        to,
		Amount:    50,
		Reference: "POST-STOP",
		Done:      make(chan error, 1),
	}
	err := proc.Submit(job)
	if err == nil {
		t.Error("expected submit to fail after Stop(), but got no error")
	}

	// Ensure the 2 original jobs completed
	mu.Lock()
	defer mu.Unlock()
	if len(executed) != 2 {
		t.Errorf("expected 2 jobs executed before stop, got %d", len(executed))
	}
}
