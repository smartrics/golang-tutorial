package async

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartrics/golang-tutorial/internal/bank"
)

/**
The processor is a responsible to handle incoming request of transfer concurrently, in a scalable and safe way.
*/

type TransferJob struct {
	From      bank.BankAccount
	To        bank.BankAccount
	Amount    float64
	Reference string
	Done      chan error // or callback
}

type Processor interface {
	Start(ctx context.Context) error
	Submit(job TransferJob) error
	Stop() error
}

type asyncProcessor struct {
	pipeline bank.TransferFunc
	workers  int
	jobs     chan TransferJob
	wg       sync.WaitGroup
	stopOnce sync.Once
	stopCh   chan struct{}
}

func NewProcessor(pipeline bank.TransferFunc, workers int) Processor {
	return &asyncProcessor{
		pipeline: pipeline,
		workers:  workers,
		jobs:     make(chan TransferJob, 100), // buffer can be tuned
		stopCh:   make(chan struct{}),
	}
}

func (p *asyncProcessor) Start(ctx context.Context) error {
	for range p.workers {
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

func (p *asyncProcessor) Stop() error {
	p.stopOnce.Do(func() {
		close(p.stopCh)
		p.wg.Wait()
	})
	return nil
}

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

func (p *asyncProcessor) process(job TransferJob) error {
	_, _, err := p.pipeline(job.From, job.To, job.Amount, job.Reference)
	return err
}
