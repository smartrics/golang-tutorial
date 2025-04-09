package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/smartrics/golang-tutorial/internal/bank"
	"github.com/smartrics/golang-tutorial/internal/bank/async"
	"github.com/smartrics/golang-tutorial/internal/engine"
	"github.com/smartrics/golang-tutorial/internal/httpapi"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Setup core
	reg := engine.NewAccountRegistry()
	svc := async.NewConcurrentBankService()
	proc := async.NewProcessor(bank.CoreTransfer(svc), 4)
	if err := proc.Start(ctx); err != nil {
		log.Fatalf("processor failed to start: %v", err)
	}

	eng := engine.NewEngine(reg, proc, svc)
	handler := httpapi.NewHandler(eng)
	router := httpapi.NewRouter(handler)

	addr := ":8080"
	log.Printf("🚀 HTTP API starting on %s", addr)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Shutdown handling
	go func() {
		<-ctx.Done()
		log.Printf("🔌 Shutting down HTTP API")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("💥 HTTP server error: %v", err)
	}
}
