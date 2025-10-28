package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func main() {
	dsn := fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable",
		getenv("DB_USER", "root"),
		getenv("DB_HOST", "cockroach"),
		getenv("DB_PORT", "26257"),
		getenv("DB_NAME", "singlespark"),
	)
	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/ping", pingHandler)
	log.Println("backend listening :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var n int
	if err := db.QueryRowContext(ctx, "SELECT 1").Scan(&n); err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("db pong"))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
