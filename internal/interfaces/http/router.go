package httpapi

import (
    "net/http"

    "github.com/go-chi/chi/v5"
)

func NewRouter(h *Handlers) http.Handler {
    r := chi.NewRouter()

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })

    r.Get("/stocks", h.ListStocks)
    r.Post("/ingest", h.Ingest)

    return r
}