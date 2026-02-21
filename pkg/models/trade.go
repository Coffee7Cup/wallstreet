package models

import (
	"time"
)

type Trade struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	CompanyID     int       `json:"company_id"`
	CompanyName   string    `json:"company_name"`
	CompanySymbol string    `json:"company_symbol"`
	TradeType     string    `json:"trade_type"` // 'BUY' or 'SELL'
	Quantity      int       `json:"quantity"`
	Price         float64   `json:"price"`
	Date          time.Time `json:"date"`
	Timestamp     time.Time `json:"timestamp"`
}
