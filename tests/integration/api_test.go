package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func TestRealServer_CreateAccountAndGetStatement(t *testing.T) {
	// 1. Run the server in background
	cmd := exec.Command("go", "run", "../../cmd/api/main.go")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer func() {
		_ = cmd.Process.Kill()
	}()

	// 2. Wait for it to boot
	time.Sleep(500 * time.Millisecond)

	// 3. Create account
	payload := []byte(`{"id":"CLI", "balance":123}`)
	res, err := http.Post("http://localhost:8080/accounts", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}

	// 4. Get statement
	resp, err := http.Get("http://localhost:8080/statement/CLI")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	defer resp.Body.Close()

	var txs []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&txs); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(txs) != 0 {
		t.Errorf("expected empty statement, got %v", txs)
	}
}
