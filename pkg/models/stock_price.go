package models

import (
	"time"
)

type StockPrice struct {
	ID            int       `json:"id"`
	CompanyID     int       `json:"company_id"`
	Date          time.Time `json:"date"`
	OpenPrice     float64   `json:"open_price"`
	ClosePrice    float64   `json:"close_price"`
	HighPrice     float64   `json:"high_price"`
	LowPrice      float64   `json:"low_price"`
	NoOfShares    int64     `json:"no_of_shares"`
	NoOfTrades    int       `json:"no_of_trades"`
	TotalTurnover float64   `json:"total_turnover"`
}
