package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/smartrics/golang-tutorial/internal/bank"
	"github.com/smartrics/golang-tutorial/internal/engine"
)

type Handler struct {
	engine engine.TransferEngine
}

func NewHandler(e engine.TransferEngine) *Handler {
	return &Handler{engine: e}
}

type auditEntryResponse struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    float64   `json:"amount"`
	Reference string    `json:"reference"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

type transferRequest struct {
	FromID    string  `json:"from_id"`
	ToID      string  `json:"to_id"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
}

type transactionResponse struct {
	ID        string  `json:"id"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
}

type createAccountRequest struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.engine.SubmitTransfer(req.FromID, req.ToID, req.Amount, req.Reference)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) Statement(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	txs, err := h.engine.GetStatement(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Transform to response shape
	out := make([]transactionResponse, len(txs))
	for i, tx := range txs {
		out[i] = transactionResponse{
			ID:        string(tx.ID()),
			From:      string(tx.From()),
			To:        string(tx.To()),
			Amount:    tx.Amount(),
			Reference: tx.Reference(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (h *Handler) DebugJobs(w http.ResponseWriter, r *http.Request) {
	var logCopy []engine.AuditEntry
	h.engine.AuditLog(&logCopy)

	out := make([]auditEntryResponse, len(logCopy))
	for i, entry := range logCopy {
		out[i] = auditEntryResponse(entry)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		http.Error(w, "account ID required", http.StatusBadRequest)
		return
	}

	// Check if account exists
	if _, err := h.engine.GetStatement(req.ID); err == nil {
		http.Error(w, "account already exists", http.StatusConflict)
		return
	}

	acc := bank.NewBankAccount(bank.AccountID(req.ID), req.Balance)
	h.engine.RegisterAccount(acc)

	w.WriteHeader(http.StatusCreated)
}
