package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type StocksRepo struct {
	DB *sql.DB
}

func NewStocksRepo(db *sql.DB) *StocksRepo {
	return &StocksRepo{DB: db}
}

// InsertMany inserta múltiples items en una única transacción.
func (r *StocksRepo) InsertMany(ctx context.Context, items []map[string]interface{}) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, raw_json)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, it := range items {
		// Mapeo defensivo: intenta extraer por campos comunes (case-insensitive aproximado)
		ticker := fmt.Sprint(it["Ticker"])
		if ticker == "" {
			ticker = fmt.Sprint(it["ticker"])
		}
		company := fmt.Sprint(it["Company"])
		if company == "" {
			company = fmt.Sprint(it["company"])
		}
		brokerage := fmt.Sprint(it["Brokerage"])
		action := fmt.Sprint(it["Action"])
		ratingFrom := fmt.Sprint(it["Rating From"])
		if ratingFrom == "" {
			ratingFrom = fmt.Sprint(it["rating_from"])
		}
		ratingTo := fmt.Sprint(it["Rating To"])
		targetFrom := fmt.Sprint(it["Target From"])
		targetTo := fmt.Sprint(it["Target To"])
		raw, _ := json.Marshal(it)

		if _, err := stmt.ExecContext(ctx, ticker, company, brokerage, action, ratingFrom, ratingTo, targetFrom, targetTo, raw); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
