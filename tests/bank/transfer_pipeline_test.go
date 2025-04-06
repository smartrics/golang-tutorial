package bank_test

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/smartrics/golang-tutorial/internal/bank"
)

func TestWithValidation(t *testing.T) {
	dummy := func(from, to bank.BankAccount, amount float64, ref string) (bank.BankAccount, bank.BankAccount, error) {
		return from, to, nil
	}
	valid := bank.WithValidation(dummy)

	acc := bank.NewBankAccount("A", 100)
	_, _, err := valid(acc, acc, 10, "ref")
	if err == nil || !strings.Contains(err.Error(), "same account") {
		t.Errorf("expected same account error, got %v", err)
	}

	_, _, err = valid(nil, acc, 10, "ref")
	if err == nil || !strings.Contains(err.Error(), "must not be nil") {
		t.Errorf("expected nil account error, got %v", err)
	}

	_, _, err = valid(acc, bank.NewBankAccount("B", 100), -1, "ref")
	if err == nil || !strings.Contains(err.Error(), "must be >= 0") {
		t.Errorf("expected negative amount error, got %v", err)
	}
}

func TestWithAudit(t *testing.T) {
	var audit []bank.AuditEntry

	auditFn := func(e bank.AuditEntry) {
		audit = append(audit, e)
	}

	dummy := func(from, to bank.BankAccount, amount float64, ref string) (bank.BankAccount, bank.BankAccount, error) {
		return from, to, nil
	}

	wrapped := bank.WithAudit(auditFn)(dummy)

	accA := bank.NewBankAccount("A", 100)
	accB := bank.NewBankAccount("B", 200)

	_, _, err := wrapped(accA, accB, 10, "PAY001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(audit) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(audit))
	}

	if audit[0].Reference != "PAY001" || !audit[0].Success {
		t.Errorf("unexpected audit contents: %+v", audit[0])
	}
}

func TestWithLoggingFailure(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil) // Reset after test

	base := func(from, to bank.BankAccount, amount float64, ref string) (bank.BankAccount, bank.BankAccount, error) {
		time.Sleep(10 * time.Millisecond)
		return nil, nil, fmt.Errorf("errorABC")
	}

	accA := bank.NewBankAccount("A", 100)
	accB := bank.NewBankAccount("B", 100)

	logged := bank.WithLogging(base)

	_, _, err := logged(accA, accB, 500, "ref")
	if err == nil {
		t.Fatalf("error expected")
	}

	logs := buf.String()
	if !strings.Contains(logs, "Initiating") || !strings.Contains(logs, "errorABC") {
		t.Errorf("expected logs to contain start and success, got:\n%s", logs)
	}
}

func TestWithLoggingSuccess(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil) // Reset after test

	base := func(from, to bank.BankAccount, amount float64, ref string) (bank.BankAccount, bank.BankAccount, error) {
		time.Sleep(10 * time.Millisecond)
		return from, to, nil
	}

	accA := bank.NewBankAccount("A", 100)
	accB := bank.NewBankAccount("B", 100)

	logged := bank.WithLogging(base)

	_, _, err := logged(accA, accB, 50, "ref")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logs := buf.String()
	if !strings.Contains(logs, "Initiating") || !strings.Contains(logs, "SUCCESS") {
		t.Errorf("expected logs to contain start and success, got:\n%s", logs)
	}
}
