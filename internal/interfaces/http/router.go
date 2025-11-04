package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NewRouter configura y devuelve el router principal del API HTTP.
func NewRouter(h *Handlers) http.Handler {
	r := chi.NewRouter()

	// Endpoint básico de salud
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Rutas de la API
	r.Route("/api", func(r chi.Router) {
		r.Get("/stocks", h.ListStocks)
		r.Post("/ingest", h.Ingest)
	})

	return r
}
