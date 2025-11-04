package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/db"
	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/usecase/stockusecase"
)

// Handlers agrupa las dependencias de la capa HTTP.
type Handlers struct {
	Repo     *db.PGRepo
	Ingestor stockusecase.Ingestor
}

// ListStocks maneja la obtención paginada de registros.
func (h *Handlers) ListStocks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 50
	offset := 0

	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			offset = n
		}
	}

	list, err := h.Repo.List(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": list,
		"meta": map[string]int{"count": len(list)},
	})
}

// Ingest ejecuta la ingesta manual vía HTTP.
func (h *Handlers) Ingest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	count, err := h.Ingestor.Ingest(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ingested": count,
	})
}
