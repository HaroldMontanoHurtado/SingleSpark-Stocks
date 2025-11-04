package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/ingestor"
	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/repository"
)

// IngestHandler maneja la ingestión manual via HTTP.
type IngestHandler struct {
	Ingestor *ingestor.Ingestor
	Repo     *repository.StocksRepo
}

// NewIngestHandler crea el handler
func NewIngestHandler(ing *ingestor.Ingestor, db *sql.DB) *IngestHandler {
	return &IngestHandler{
		Ingestor: ing,
		Repo:     repository.NewStocksRepo(db),
	}
}

// ServeHTTP fuerza la ingestión (fetch -> persist). Retorna 202 si se aceptó.
func (h *IngestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	// simple: consumir páginas hasta que next == ""
	next := ""
	all := make([]map[string]interface{}, 0)
	for {
		items, nextPage, err := h.Ingestor.FetchPage(ctx, next)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(items) == 0 {
			break
		}
		// convertir items (type Item) a []map[string]interface{}
		for _, it := range items {
			all = append(all, it)
		}
		if nextPage == "" {
			break
		}
		next = nextPage
	}

	if len(all) == 0 {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("no data"))
		return
	}

	if err := h.Repo.InsertMany(ctx, all); err != nil {
		http.Error(w, "insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	b, _ := json.Marshal(map[string]interface{}{"inserted": len(all)})
	w.Write(b)
}
