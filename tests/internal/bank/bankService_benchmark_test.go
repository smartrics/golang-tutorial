package bank_test

import (
	"testing"

	"github.com/smartrics/golang-tutorial/internal/bank"
)

func BenchmarkBankService_Transfer2(b *testing.B) {
	b.Logf("Benchmark running")
	bs := bank.NewBankService()
	for b.Loop() {
		from := bank.NewBankAccount("from", 1_000_000_000)
		to := bank.NewBankAccount("to", 0)
		_, _, err := bs.Transfer(from, to, 100, "bench")
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
