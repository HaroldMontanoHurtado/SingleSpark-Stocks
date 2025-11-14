package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/dominio/stock"
)

// PGRepo implementa Repository usando pgxpool.
type PGRepo struct {
	pool *pgxpool.Pool
}

// NewPGRepo conecta y devuelve el repo.
func NewPGRepo(ctx context.Context, connString string) (*PGRepo, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	// Opciones tuning si se desee
	cfg.MaxConns = 8

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &PGRepo{pool: pool}, nil
}

func (p *PGRepo) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// Save inserta o actualiza un stock (upsert básico).
func (p *PGRepo) Save(s *stock.Stock) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	var rawJSON json.RawMessage
	if s.Raw == "" {
		rawJSON = json.RawMessage([]byte("null"))
	} else {
		rawJSON = json.RawMessage([]byte(s.Raw))
	}
	_, err := p.pool.Exec(context.Background(),
		`INSERT INTO stocks (id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, raw, created_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
        ON CONFLICT (id) DO UPDATE
        SET ticker = EXCLUDED.ticker,
            company = EXCLUDED.company,
            brokerage = EXCLUDED.brokerage,
            action = EXCLUDED.action,
            rating_from = EXCLUDED.rating_from,
            rating_to = EXCLUDED.rating_to,
            target_from = EXCLUDED.target_from,
            target_to = EXCLUDED.target_to,
            raw = EXCLUDED.raw,
            created_at = EXCLUDED.created_at`,
		s.ID, s.Ticker, s.Company, s.Brokerage, s.Action, s.RatingFrom, s.RatingTo, s.TargetFrom, s.TargetTo, rawJSON, s.CreatedAt)
	return err
}

// SaveBatch inserta muchos registros en una transacción.
func (p *PGRepo) SaveBatch(stocks []*stock.Stock) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	stmt := `INSERT INTO stocks (id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, raw, created_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	    ON CONFLICT (id) DO UPDATE
        SET ticker = EXCLUDED.ticker,
	        company = EXCLUDED.company,
	        brokerage = EXCLUDED.brokerage,
	        action = EXCLUDED.action,
	        rating_from = EXCLUDED.rating_from,
	        rating_to = EXCLUDED.rating_to,
	        target_from = EXCLUDED.target_from,
	        target_to = EXCLUDED.target_to,
	        raw = EXCLUDED.raw,
	        created_at = EXCLUDED.created_at`

	for _, s := range stocks {
		if s == nil {
			continue
		}
		if s.ID == "" {
			s.ID = uuid.New().String()
		}
		var rawJSON json.RawMessage
		if s.Raw == "" {
			rawJSON = json.RawMessage([]byte("null"))
		} else {
			rawJSON = json.RawMessage([]byte(s.Raw))
		}
		if _, err := tx.Exec(ctx, stmt,
			s.ID, s.Ticker, s.Company, s.Brokerage, s.Action, s.RatingFrom, s.RatingTo, s.TargetFrom, s.TargetTo, rawJSON, s.CreatedAt); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// List retorna con paginación.
func (p *PGRepo) List(limit, offset int) ([]*stock.Stock, error) {
	rows, err := p.pool.Query(context.Background(),
		`SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, raw, created_at
        FROM stocks
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*stock.Stock
	for rows.Next() {
		var s stock.Stock
		var raw []byte
		if err := rows.Scan(&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action, &s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &raw, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.Raw = string(raw)
		out = append(out, &s)
	}
	return out, nil
}

// FindByTicker devuelve items por ticker.
func (p *PGRepo) FindByTicker(ticker string) ([]*stock.Stock, error) {
	rows, err := p.pool.Query(context.Background(),
		`SELECT id, ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, raw, created_at
        FROM stocks
        WHERE ticker = $1
        ORDER BY created_at DESC`, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*stock.Stock
	for rows.Next() {
		var s stock.Stock
		var raw []byte
		if err := rows.Scan(&s.ID, &s.Ticker, &s.Company, &s.Brokerage, &s.Action, &s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &raw, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.Raw = string(raw)
		out = append(out, &s)
	}
	return out, nil
}

// ParseExternalItem intenta mapear un item externo (map[string]interface{}) a *stock.Stock
// Esta función la exportamos para usarla desde el usecase/ingestor.
func ParseExternalItem(m map[string]interface{}) (*stock.Stock, error) {
	// Aquí hay heurística: adapta esto al JSON real de la API.
	getStr := func(keys ...string) string {
		for _, k := range keys {
			if v, ok := m[k]; ok && v != nil {
				if s, ok := v.(string); ok {
					return s
				}
				if b, err := json.Marshal(v); err == nil {
					return string(b)
				}
			}
		}
		return ""
	}

	id := getStr("id", "uuid", "ID")
	ticker := getStr("ticker", "symbol", "symbol_ticker")
	company := getStr("company", "company_name", "title", "name")
	brokerage := getStr("brokerage", "source")
	action := getStr("action")
	ratingFrom := getStr("rating_from", "ratingFrom")
	ratingTo := getStr("rating_to", "ratingTo")
	targetFrom := getStr("target_from", "targetFrom")
	targetTo := getStr("target_to", "targetTo")

	rawb, _ := json.Marshal(m)

	if ticker == "" && company == "" {
		return nil, fmt.Errorf("missing ticker/company in external item")
	}

	return &stock.Stock{
		ID:         id,
		Ticker:     ticker,
		Company:    company,
		Brokerage:  brokerage,
		Action:     action,
		RatingFrom: ratingFrom,
		RatingTo:   ratingTo,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Raw:        string(rawb),
		CreatedAt:  time.Now().UTC(),
	}, nil
}
