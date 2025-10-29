package httpapi

import (
    //"context"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/db"
    stockusecase "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/usecase/stock"
)

type Handlers struct {
    Repo     *db.PGRepo
    Ingestor stockusecase.Ingestor
}

func (h *Handlers) ListStocks(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    l := 50
    o := 0
    if v := q.Get("limit"); v != "" {
        if n, err := strconv.Atoi(v); err == nil {
            l = n
        }
    }
    if v := q.Get("offset"); v != "" {
        if n, err := strconv.Atoi(v); err == nil {
            o = n
        }
    }
    list, err := h.Repo.List(l, o)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "data": list,
        "meta": map[string]int{"count": len(list)},
    })
}

func (h *Handlers) Ingest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    n, err := h.Ingestor.IngestOnce(ctx)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "ingested": n,
    })
}
