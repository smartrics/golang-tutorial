package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/smartrics/golang-tutorial/internal/bank"
	"github.com/smartrics/golang-tutorial/internal/bank/async"
	"github.com/smartrics/golang-tutorial/internal/engine"
	"github.com/smartrics/golang-tutorial/internal/httpapi"
)

func TestTransferE2E(t *testing.T) {
	// Setup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reg := engine.NewAccountRegistry()
	from := bank.NewBankAccount("X", 1000)
	to := bank.NewBankAccount("Y", 0)
	reg.Register(from)
	reg.Register(to)

	svc := async.NewConcurrentBankService()
	proc := async.NewProcessor(bank.CoreTransfer(svc), 1)
	_ = proc.Start(ctx)

	eng := engine.NewEngine(reg, proc, svc)
	handler := httpapi.NewHandler(eng)
	router := httpapi.NewRouter(handler)

	// Simulate POST /transfer
	payload := map[string]interface{}{
		"from_id":   "X",
		"to_id":     "Y",
		"amount":    150.0,
		"reference": "e2e-test",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusAccepted {
		t.Fatalf("unexpected status code: %d", res.Code)
	}

	// Wait for job to process
	time.Sleep(300 * time.Millisecond)

	// Simulate GET /statement/X
	req = httptest.NewRequest(http.MethodGet, "/statement/X", nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected OK, got %d", res.Code)
	}

	var txs []bank.Transaction
	if err := json.NewDecoder(res.Body).Decode(&txs); err != nil {
		t.Fatalf("failed to decode statement: %v", err)
	}

	if len(txs) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(txs))
	}
}

func TestCreateAccountAndTransferE2E(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reg := engine.NewAccountRegistry()
	svc := async.NewConcurrentBankService()
	proc := async.NewProcessor(bank.CoreTransfer(svc), 1)
	_ = proc.Start(ctx)

	eng := engine.NewEngine(reg, proc, svc)
	handler := httpapi.NewHandler(eng)
	router := httpapi.NewRouter(handler)

	// Create "A"
	body := []byte(`{"id":"A","balance":500}`)
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.Code)
	}

	// Create "B"
	body = []byte(`{"id":"B","balance":0}`)
	req = httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.Code)
	}

	// Transfer
	body = []byte(`{"from_id":"A","to_id":"B","amount":100,"reference":"http-test"}`)
	req = httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", res.Code)
	}

	time.Sleep(300 * time.Millisecond)

	// Check B's statement
	req = httptest.NewRequest(http.MethodGet, "/statement/B", nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	var txs []bank.Transaction
	if err := json.NewDecoder(res.Body).Decode(&txs); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(txs) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(txs))
	}
}

func TestCreateAccountAndGetEmptyStatement(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup engine
	reg := engine.NewAccountRegistry()
	svc := async.NewConcurrentBankService()
	proc := async.NewProcessor(bank.CoreTransfer(svc), 1)
	_ = proc.Start(ctx)

	eng := engine.NewEngine(reg, proc, svc)
	handler := httpapi.NewHandler(eng)
	router := httpapi.NewRouter(handler)

	// 1. Create account "Z"
	body := []byte(`{"id":"Z","balance":1234}`)
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", res.Code)
	}

	// 2. Retrieve statement for "Z"
	req = httptest.NewRequest(http.MethodGet, "/statement/Z", nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", res.Code)
	}

	var txs []bank.Transaction
	if err := json.NewDecoder(res.Body).Decode(&txs); err != nil {
		t.Fatalf("error decoding statement: %v", err)
	}
	if len(txs) != 0 {
		t.Errorf("expected empty statement, got %v", txs)
	}
}
