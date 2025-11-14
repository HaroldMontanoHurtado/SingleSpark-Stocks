package httpapi

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/db"
    "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/usecase/stockusecase"
)

type Handlers struct {
    Repo     *db.PGRepo
    Ingestor stockusecase.Ingestor
}

// GET /api/stocks
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

// POST /api/ingest
func (h *Handlers) Ingest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    ctx := r.Context()
    count, err := h.Ingestor.Ingest(ctx)
    if err != nil {
        log.Println("ingest error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "ingested": count,
    })
}
