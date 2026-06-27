package main

import (
	"complaint-portal/handlers"
	"complaint-portal/store"
	"fmt"
	"log"
	"net/http"
)

const (
	port        = ":8080"
	adminSecret = "ADMIN-SECRET-2024"
)

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

func main() {
	s := store.New(adminSecret)
	h := &handlers.Handler{Store: s}

	mux := http.NewServeMux()

	mux.HandleFunc("/register", loggingMiddleware(h.Register))
	mux.HandleFunc("/login", loggingMiddleware(h.Login))
	mux.HandleFunc("/submitComplaint", loggingMiddleware(h.SubmitComplaint))
	mux.HandleFunc("/getAllComplaintsForUser", loggingMiddleware(h.GetAllComplaintsForUser))
	mux.HandleFunc("/getAllComplaintsForAdmin", loggingMiddleware(h.GetAllComplaintsForAdmin))
	mux.HandleFunc("/viewComplaint", loggingMiddleware(h.ViewComplaint))
	mux.HandleFunc("/resolveComplaint", loggingMiddleware(h.ResolveComplaint))

	fmt.Printf("Complaint Portal API running on http://localhost%s\n", port)
	fmt.Printf("Admin Secret Code: %s\n", adminSecret)
	fmt.Println("─────────────────────────────────────────")

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
