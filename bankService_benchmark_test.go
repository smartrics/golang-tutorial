package main

import "testing"

func BenchmarkBankService_Transfer2(b *testing.B) {
	b.Logf("Benchmark running")
	bs := NewBankService()
	for b.Loop() {
		from := NewBankAccount("from", 1_000_000_000)
		to := NewBankAccount("to", 0)
		_, _, err := bs.Transfer(from, to, 100, "bench")
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
