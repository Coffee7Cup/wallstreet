package models

import (
	"time"
)

type Trade struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CompanyID int       `json:"company_id"`
	TradeType string    `json:"trade_type"` // 'BUY' or 'SELL'
	Quantity  int       `json:"quantity"`
	Date      time.Time `json:"date"`
	Timestamp time.Time `json:"timestamp"`
}
