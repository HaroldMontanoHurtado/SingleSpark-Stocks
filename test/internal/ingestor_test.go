package internal_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HaroldMontanoHurtado/SingleSpark-Stocks/internal/infrastructure/ingestor"
)

// TestFetchPage simula una respuesta HTTP para probar el ingestor.
func TestFetchPage(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		resp := map[string]interface{}{
			"data": []map[string]interface{}{
				{"symbol": "AAPL", "price": 170.5},
				{"symbol": "GOOGL", "price": 135.2},
			},
			"next": "",
			"page": page,
		}
		b, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}))
	defer mockServer.Close()

	// Crear el ingestor con la función real
	ing := ingestor.New(mockServer.URL, "fake_token")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, next, err := ing.FetchPage(ctx, "")
	if err != nil {
		t.Fatalf("FetchPage failed: %v", err)
	}

	if len(items) == 0 {
		t.Fatalf("Expected non-empty items, got %d", len(items))
	}

	if next != "" {
		t.Fatalf("Expected next='' but got '%s'", next)
	}

	if items[0]["symbol"] != "AAPL" {
		t.Errorf("Expected symbol=AAPL, got %v", items[0]["symbol"])
	}
}

// TestFetchError verifica el manejo de errores HTTP.
func TestFetchError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	ing := ingestor.New(mockServer.URL, "fake_token")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, _, err := ing.FetchPage(ctx, "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
