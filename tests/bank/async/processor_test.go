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
