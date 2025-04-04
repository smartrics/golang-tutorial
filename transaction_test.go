package main

import (
	"testing"
)

func TestNewTransaction(t *testing.T) {
	tx := NewTransaction("123", "Alice", "Bob", 100.0, "test-ref")
	if tx == nil {
		t.Fatal("Expected a new Transaction, got nil")
	}
	if tx.ID() != "123" {
		t.Errorf("Expected ID '123', got '%s'", tx.ID())
	}
	if tx.From() != "Alice" {
		t.Errorf("Expected From 'Alice', got '%s'", tx.From())
	}
	if tx.To() != "Bob" {
		t.Errorf("Expected To 'Bob', got '%s'", tx.To())
	}
	if tx.Reference() != "test-ref" {
		t.Errorf("Expected Reference 'test-ref', got '%s'", tx.To())
	}
}
