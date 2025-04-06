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

func TestTransferFunc_PipelineExecution_Success(t *testing.T) {
	service := bank.NewBankService()

	var audit []bank.AuditEntry

	auditFn := func(auditEntry bank.AuditEntry) {
		audit = append(audit, auditEntry)
	}

	transfer := bank.WithLogging(
		bank.WithAudit(auditFn)(
			bank.CoreTransfer(service),
		),
	)

	from := bank.NewBankAccount("FROM", 1000)
	to := bank.NewBankAccount("TO", 500)

	newFrom, newTo, err := transfer(from, to, 250, "INV-999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate balances
	if newFrom.Balance() != 750 {
		t.Errorf("expected newFrom balance 750, got %.2f", newFrom.Balance())
	}
	if newTo.Balance() != 750 {
		t.Errorf("expected newTo balance 750, got %.2f", newTo.Balance())
	}

	// Validate audit log
	if len(audit) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(audit))
	}
	entry := audit[0]
	if !entry.Success || entry.Reference != "INV-999" {
		t.Errorf("unexpected audit entry: %+v", entry)
	}

	// Validate transaction was recorded
	stmtFrom, _ := service.GetStatement(newFrom)
	stmtTo, _ := service.GetStatement(newTo)

	if len(stmtFrom) == 0 || len(stmtTo) == 0 {
		t.Error("expected statements to include transaction")
	}

}

func TestTransferFunc_PipelineExecution_Failure(t *testing.T) {
	service := bank.NewBankService()

	var audit []bank.AuditEntry

	auditFn := func(auditEntry bank.AuditEntry) {
		audit = append(audit, auditEntry)
	}

	transfer := bank.WithAudit(auditFn)(
		bank.CoreTransfer(service),
	)

	from := bank.NewBankAccount("A1", 50)
	to := bank.NewBankAccount("A2", 0)

	_, _, err := transfer(from, to, 100, "failRef")

	if err == nil {
		t.Fatalf("expected an error")
	}

	// Validate audit log
	if len(audit) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(audit))
	}
	entry := audit[0]
	if entry.Success || entry.Reference != "failRef" {
		t.Errorf("unexpected audit entry: %+v", entry)
	}

}

func TestPipeline_WithAccountWithCounter(t *testing.T) {
	service := bank.NewBankService()

	var audit []bank.AuditEntry
	auditFn := func(e bank.AuditEntry) {
		audit = append(audit, e)
	}

	transfer := bank.WithLogging(
		bank.WithValidation(
			bank.WithAudit(auditFn)(
				bank.CoreTransfer(service),
			),
		),
	)

	baseFrom := bank.NewBankAccount("FROM", 1000)
	baseTo := bank.NewBankAccount("TO", 500)

	from := bank.NewAccountWithCounter(baseFrom)
	to := bank.NewAccountWithCounter(baseTo)

	resFrom, resTo, err := transfer(from, to, 200, "DEMO-100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Type assert back to AccountWithCounter to inspect counters
	fromWithCount, ok := resFrom.(*bank.AccountWithCounter)
	if !ok {
		t.Fatalf("expected AccountWithCounter, got %T", resFrom)
	}
	toWithCount, ok := resTo.(*bank.AccountWithCounter)
	if !ok {
		t.Fatalf("expected AccountWithCounter, got %T", resTo)
	}

	if fromWithCount.Count() != 1 {
		t.Errorf("expected from count 1, got %d", fromWithCount.Count())
	}
	if toWithCount.Count() != 1 {
		t.Errorf("expected to count 1, got %d", toWithCount.Count())
	}

	// Bonus: assert balance updated correctly
	if fromWithCount.Balance() != 800 || toWithCount.Balance() != 700 {
		t.Errorf("unexpected balances: from=%.2f, to=%.2f", fromWithCount.Balance(), toWithCount.Balance())
	}

	// Bonus: audit and statement check
	if len(audit) != 1 {
		t.Errorf("expected 1 audit entry, got %d", len(audit))
	}
	stmt, _ := service.GetStatement(fromWithCount)
	if len(stmt) == 0 {
		t.Errorf("expected statement for from account")
	}
}
