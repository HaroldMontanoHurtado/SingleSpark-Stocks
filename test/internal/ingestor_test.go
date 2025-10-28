package internal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"context"
	"io"
	"strings"

	"github.com/your/module/internal/infrastructure/ingestor"
)

func TestFetchPage_Simple(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"data":[{"Ticker":"ABC","Company":"ACME Inc","Action":"reiterated by","Rating From":"Buy","Rating To":"Buy","Target From":"10","Target To":"12"}]}`)
	}))
	defer ts.Close()

	ing := ingestor.New(ts.URL, "dummy")
	items, next, err := ing.FetchPage(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item got %d", len(items))
	}
	if next != "" {
		t.Fatalf("expected no next, got %s", next)
	}
	if items[0]["Ticker"] != "ABC" && items[0]["ticker"] != "ABC" {
		// ok: check keys may be capitalized
		t.Fatalf("unexpected ticker: %#v", items[0])
	}
}
