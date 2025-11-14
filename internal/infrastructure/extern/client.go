package extern

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client es un cliente simple para la API externa.
type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

// NewClient crea el cliente.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTP: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// FetchList obtiene la lista desde la API externa.
// Devuelve: items (slice de map[string]interface{}), raw bytes, error.
func (c *Client) FetchList(ctx context.Context, query string) ([]map[string]interface{}, []byte, error) {
	// Si la API admite query params, podríamos concatenar query; por ahora usamos BaseURL directo.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL, nil)
	if err != nil {
		return nil, nil, err
	}

	// Encabezado de autenticación (si la API lo requiere)
	if c.APIKey != "" {
		// Ajusta si la API usa Authorization Bearer o x-api-key
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
		req.Header.Set("X-API-Key", c.APIKey) // redundante; algunos endpoints usan uno u otro
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, body, fmt.Errorf("external api status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Intentamos decodificar a un formato flexible: una lista de objetos.
	var parsed interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, body, err
	}

	// Normalizar a []map[string]interface{}
	var items []map[string]interface{}
	switch v := parsed.(type) {
	case []interface{}:
		for _, it := range v {
			if m, ok := it.(map[string]interface{}); ok {
				items = append(items, m)
			}
		}
	case map[string]interface{}:
		// Muchas APIs devuelven {"data": [...]} o {"items": [...]}
		if maybe, ok := v["data"]; ok {
			if arr, ok := maybe.([]interface{}); ok {
				for _, it := range arr {
					if m, ok := it.(map[string]interface{}); ok {
						items = append(items, m)
					}
				}
			}
		} else if maybe, ok := v["items"]; ok {
			if arr, ok := maybe.([]interface{}); ok {
				for _, it := range arr {
					if m, ok := it.(map[string]interface{}); ok {
						items = append(items, m)
					}
				}
			}
		} else {
			// si el map es un solo item, devolverlo como slice de 1
			items = append(items, v)
		}
	default:
		// no reconocido
		return nil, body, fmt.Errorf("unexpected external payload type %T", v)
	}

	return items, body, nil
}
