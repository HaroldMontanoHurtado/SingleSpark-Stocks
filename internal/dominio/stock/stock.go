package stock

import "time"

type Stock struct {
    ID         string    `json:"id"`
    Ticker     string    `json:"ticker"`
    Company    string    `json:"company,omitempty"`
    Brokerage  string    `json:"brokerage,omitempty"`
    Action     string    `json:"action,omitempty"`
    RatingFrom string    `json:"rating_from,omitempty"`
    RatingTo   string    `json:"rating_to,omitempty"`
    TargetFrom string    `json:"target_from,omitempty"`
    TargetTo   string    `json:"target_to,omitempty"`
    Raw        string    `json:"raw,omitempty"` // raw JSON string
    CreatedAt  time.Time `json:"created_at"`
}
