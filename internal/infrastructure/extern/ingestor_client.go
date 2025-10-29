package extern

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type Client struct {
    url   string
    token string
    http  *http.Client
}

func NewClient(url, token string) *Client {
    return &Client{
        url:   url,
        token: token,
        http: &http.Client{
            Timeout: 15 * time.Second,
        },
    }
}

// FetchList performs GET and returns []map[string]interface{} (raw items)
func (c *Client) FetchList(ctx context.Context, nextPage string) ([]map[string]interface{}, string, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
    if err != nil {
        return nil, "", err
    }
    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }
    req.Header.Set("Content-Type", "application/json")

    q := req.URL.Query()
    if nextPage != "" {
        q.Set("next_page", nextPage)
        req.URL.RawQuery = q.Encode()
    }

    resp, err := c.http.Do(req)
    if err != nil {
        return nil, "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        body, _ := io.ReadAll(resp.Body)
        return nil, "", fmt.Errorf("external api returned %d: %s", resp.StatusCode, string(body))
    }

    var data map[string]interface{}
    dec := json.NewDecoder(resp.Body)
    dec.UseNumber()
    if err := dec.Decode(&data); err != nil {
        return nil, "", err
    }

    // The response shape may vary; try to extract a list
    if list, ok := data["data"].([]interface{}); ok {
        out := make([]map[string]interface{}, 0, len(list))
        for _, item := range list {
            if m, ok := item.(map[string]interface{}); ok {
                out = append(out, m)
            }
        }
        // check for next_page token
        next := ""
        if np, ok := data["next_page"].(string); ok {
            next = np
        }
        return out, next, nil
    }

    // fallback: maybe the top-level is an array
    if arr, ok := data["items"].([]interface{}); ok {
        out := make([]map[string]interface{}, 0, len(arr))
        for _, item := range arr {
            if m, ok := item.(map[string]interface{}); ok {
                out = append(out, m)
            }
        }
        return out, "", nil
    }

    // if root is array
    if arrRoot, ok := data[""].([]interface{}); ok {
        _ = arrRoot
    }

    // as last fallback, try to see if top-level is an array by re-decoding raw body (we already consumed it)
    // but since we already decoded, try to convert whole map to a single-item list
    return []map[string]interface{}{data}, "", nil
}
