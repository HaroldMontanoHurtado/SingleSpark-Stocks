package httpapi

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handlers) http.Handler {
    r := chi.NewRouter()

    // Middlewares generales
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // CORS básico
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
            
            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    })

    // Health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok"}`))
    })

    // API
    r.Route("/api", func(r chi.Router) {
        r.Get("/stocks", h.ListStocks)
        r.Post("/ingest", h.Ingest)
    })

    return r
}
