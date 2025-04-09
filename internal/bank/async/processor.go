package async

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartrics/golang-tutorial/internal/bank"
)

// TransferJob represents a single banking transfer task.
//
// It contains:
// - The source and destination BankAccount
// - The amount to transfer and a reference string
// - An optional Done channel to signal completion and return any error
// - An optional context (Ctx) to enforce timeout or cancellation per job
type TransferJob struct {
	From      bank.BankAccount
	To        bank.BankAccount
	Amount    float64
	Reference string
	Done      chan error
	Ctx       context.Context
}

// Processor defines the interface for a concurrent transfer processor.
//
// It can be started and stopped gracefully and accepts incoming transfer jobs
// via Submit(). The processor runs a configurable number of worker goroutines
// that pull jobs from a channel and execute them using a TransferFunc pipeline.
type Processor interface {
	Start(ctx context.Context) error
	Submit(job TransferJob) error
	Stop() error
}

// asyncProcessor is a concrete implementation of Processor.
//
// It executes TransferJobs using a user-provided bank.TransferFunc pipeline,
// running multiple concurrent workers and supporting graceful shutdown,
// per-job context cancellation, and error reporting via channels.
type asyncProcessor struct {
	pipeline bank.TransferFunc // Decorated transfer logic (e.g., logging, audit, validation)
	workers  int               // Number of concurrent workers
	jobs     chan TransferJob  // Buffered channel of submitted jobs
	wg       sync.WaitGroup    // Tracks live workers for graceful shutdown
	stopOnce sync.Once         // Ensures Stop() is only called once
	stopCh   chan struct{}     // Signals workers to shut down
}

// NewProcessor creates a new asynchronous processor with the given pipeline and number of workers.
//
// The returned Processor can be started using Start() and accepts jobs via Submit().
func NewProcessor(pipeline bank.TransferFunc, workers int) Processor {
	return &asyncProcessor{
		pipeline: pipeline,
		workers:  workers,
		jobs:     make(chan TransferJob, 100), // Queue size is configurable
		stopCh:   make(chan struct{}),
	}
}

// Start launches the configured number of worker goroutines.
//
// Workers listen for incoming jobs on the internal queue and stop gracefully
// when the provided context is cancelled or Stop() is called.
func (p *asyncProcessor) Start(ctx context.Context) error {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-p.stopCh:
					return
				case job := <-p.jobs:
					err := p.process(job)
					if job.Done != nil {
						job.Done <- err
					}
				}
			}
		}()
	}
	return nil
}

// Stop signals all workers to terminate and waits for them to finish.
//
// This method is safe to call multiple times and guarantees clean shutdown.
func (p *asyncProcessor) Stop() error {
	p.stopOnce.Do(func() {
		close(p.stopCh)
		p.wg.Wait()
	})
	return nil
}

// Submit queues a transfer job for asynchronous processing.
//
// It returns an error if the processor has been stopped or if the job queue is full.
func (p *asyncProcessor) Submit(job TransferJob) error {
	select {
	case <-p.stopCh:
		return fmt.Errorf("processor stopped")
	default:
		select {
		case p.jobs <- job:
			return nil
		default:
			return fmt.Errorf("job queue full")
		}
	}
}

// process executes a single TransferJob using the configured pipeline.
//
// It respects the per-job context (if provided) and returns a timeout/cancellation
// error if the context expires before the transfer completes.
func (p *asyncProcessor) process(job TransferJob) error {
	ctx := job.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	done := make(chan struct{})
	var err error

	go func() {
		_, _, err = p.pipeline(job.From, job.To, job.Amount, job.Reference)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("job timed out or cancelled: %w", ctx.Err())
	case <-done:
		return err
	}
}
