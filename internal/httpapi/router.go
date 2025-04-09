package httpapi

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("➡️  %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func NewRouter(h *Handler) http.Handler {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/transfer", h.Transfer).Methods("POST")
	r.HandleFunc("/statement/{id}", h.Statement).Methods("GET")
	r.HandleFunc("/accounts", h.CreateAccount).Methods("POST")
	r.HandleFunc("/debug/jobs", h.DebugJobs).Methods("GET")
	return r
}
