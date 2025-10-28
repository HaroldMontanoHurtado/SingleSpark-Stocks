package ingestor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Item representa un objeto genérico que devuelve la API.
// Lo dejamos como map[string]interface{} en la persistencia para máxima flexibilidad.
type Item = map[string]interface{}

type Ingestor struct {
	HTTP     *http.Client
	BaseURL  string
	APIToken string
}

// New crea un ingestor con timeout razonable.
func New(baseURL, apiToken string) *Ingestor {
	return &Ingestor{
		HTTP: &http.Client{Timeout: 20 * time.Second},
		BaseURL: baseURL,
		APIToken: apiToken,
	}
}

// FetchPage obtiene una página y retorna items y el token/clave de next page si existe.
func (i *Ingestor) FetchPage(ctx context.Context, nextPage string) ([]Item, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, i.BaseURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+i.APIToken)
	q := req.URL.Query()
	if nextPage != "" {
		q.Set("next_page", nextPage)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := i.HTTP.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("bad status: %s", resp.Status)
	}

	var payload struct {
		Data     []Item `json:"data"`      // muchos endpoints usan "data" o "results"
		Results  []Item `json:"results"`
		NextPage string `json:"next_page"`
		Next     string `json:"next"`
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&payload); err != nil {
		// intenta decodificar en forma flexible
		var generic map[string]interface{}
		if err2 := json.NewDecoder(resp.Body).Decode(&generic); err2 != nil {
			return nil, "", fmt.Errorf("decode: %w", err)
		}
		return nil, "", nil
	}

	// Preferencia: Data > Results > fallback empty
	var items []Item
	if len(payload.Data) > 0 {
		items = payload.Data
	} else if len(payload.Results) > 0 {
		items = payload.Results
	} else {
		// intenta decodificar todo como array
		// (En caso de que la API devuelva directamente un array)
		_, _ = resp.Body.Seek(0, 0) // no siempre posible; skip
		items = []Item{}
	}

	// next token puede estar en varios campos
	next := payload.NextPage
	if next == "" {
		next = payload.Next
	}

	return items, next, nil
}
