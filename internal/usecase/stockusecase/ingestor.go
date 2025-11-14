package stockusecase

import (
	"context"
	"fmt"
	"log"

	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/dominio/stock"
	repo "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/db"
	extern "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/extern"
)

type Ingestor interface {
    Ingest(ctx context.Context) (int, error)
}

type ingestor struct {
    client *extern.Client
    repo   *repo.PGRepo
}

func NewIngestor(client *extern.Client, repo *repo.PGRepo) Ingestor {
    return &ingestor{client: client, repo: repo}
}

func (i *ingestor) Ingest(ctx context.Context) (int, error) {
    items, _, err := i.client.FetchList(ctx, "")
    if err != nil {
        return 0, err
    }
    var parsed []*stock.Stock
    for _, it := range items {
        s, perr := repo.ParseExternalItem(it)
        if perr != nil {
            log.Printf("parse error on item: %v", perr)
        }
        parsed = append(parsed, s)
    }
    if len(parsed) == 0 {
        return 0, fmt.Errorf("no items parsed")
    }
    if err := i.repo.SaveBatch(parsed); err != nil {
        return 0, err
    }
    return len(parsed), nil
}
