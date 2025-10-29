package stockusecase

import (
    "context"
    "fmt"

    extern "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/extern"
    repo "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/db"
    "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/dominio/stock"
)

type Ingestor interface {
    IngestOnce(ctx context.Context) (int, error)
}

type ingestor struct {
    client *extern.Client
    repo   *repo.PGRepo
}

func NewIngestor(client *extern.Client, repo *repo.PGRepo) Ingestor {
    return &ingestor{client: client, repo: repo}
}

func (i *ingestor) IngestOnce(ctx context.Context) (int, error) {
    items, _, err := i.client.FetchList(ctx, "")
    if err != nil {
        return 0, err
    }
    var parsed []*stock.Stock
    for _, it := range items {
        s, perr := repo.ParseExternalItem(it)
        if perr != nil {
            // continue but log
            continue
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
