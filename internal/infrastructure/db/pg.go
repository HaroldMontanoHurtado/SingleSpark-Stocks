package db

import (
    "context"
    "encoding/json"
    //"fmt"
    "time"

    "github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/dominio/stock"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PGRepo struct {
    pool *pgxpool.Pool
}

func NewPostgresRepository(connString string) (*PGRepo, error) {
    pool, err := pgxpool.New(context.Background(), connString)
    if err != nil {
        return nil, err
    }
    return &PGRepo{pool: pool}, nil
}

func (p *PGRepo) Save(s *stock.Stock) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err := p.pool.Exec(ctx, `
        INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, raw)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
    `, s.Ticker, s.Company, s.Brokerage, s.Action, s.RatingFrom, s.RatingTo, s.TargetFrom, s.TargetTo, s.Raw)
    return err
}

func (p *PGRepo) SaveBatch(stocks []*stock.Stock) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    tx, err := p.pool.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    for _, s := range stocks {
        _, err := tx.Exec(ctx, `
            INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, raw)
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        `, s.Ticker, s.Company, s.Brokerage, s.Action, s.RatingFrom, s.RatingTo, s.TargetFrom, s.TargetTo, s.Raw)
        if err != nil {
            return err
        }
    }

    return tx.Commit(ctx)
}

func (p *PGRepo) List(limit, offset int) ([]*stock.Stock, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    rows, err := p.pool.Query(ctx, `
        SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, raw, created_at
        FROM stocks ORDER BY created_at DESC LIMIT $1 OFFSET $2
    `, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []*stock.Stock
    for rows.Next() {
        var s stock.Stock
        var rawBytes []byte
        if err := rows.Scan(&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action, &s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &rawBytes, &s.CreatedAt); err != nil {
            return nil, err
        }
        s.Raw = string(rawBytes)
        out = append(out, &s)
    }
    return out, nil
}

func (p *PGRepo) FindByTicker(ticker string) ([]*stock.Stock, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    rows, err := p.pool.Query(ctx, `
        SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, raw, created_at
        FROM stocks WHERE ticker = $1 ORDER BY created_at DESC
    `, ticker)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []*stock.Stock
    for rows.Next() {
        var s stock.Stock
        var rawBytes []byte
        if err := rows.Scan(&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action, &s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &rawBytes, &s.CreatedAt); err != nil {
            return nil, err
        }
        s.Raw = string(rawBytes)
        out = append(out, &s)
    }
    return out, nil
}

// Helper: parse raw items from external API into domain Stock
func ParseExternalItem(m map[string]interface{}) (*stock.Stock, error) {
    // flexible parsing — adapt to JSON structure you receive
    b, _ := json.Marshal(m)
    s := &stock.Stock{
        Raw: string(b),
    }
    // attempt to map common fields
    if v, ok := m["Ticker"].(string); ok {
        s.Ticker = v
    }
    if v, ok := m["ticker"].(string); ok && s.Ticker == "" {
        s.Ticker = v
    }
    if v, ok := m["Company"].(string); ok {
        s.Company = v
    }
    if v, ok := m["company"].(string); ok && s.Company == "" {
        s.Company = v
    }
    if v, ok := m["Brokerage"].(string); ok {
        s.Brokerage = v
    }
    if v, ok := m["Action"].(string); ok {
        s.Action = v
    }
    if v, ok := m["Rating From"].(string); ok {
        s.RatingFrom = v
    }
    if v, ok := m["Rating To"].(string); ok {
        s.RatingTo = v
    }
    if v, ok := m["Target From"].(string); ok {
        s.TargetFrom = v
    }
    if v, ok := m["Target To"].(string); ok {
        s.TargetTo = v
    }

    // fallback: inspect common lowercased keys
    if s.Ticker == "" {
        if v, ok := m["ticker_symbol"].(string); ok {
            s.Ticker = v
        }
    }

    return s, nil
}
