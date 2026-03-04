package models

import (
	"time"
)

type Fundamental struct {
	ID              int       `json:"id"`
	CompanyID       int       `json:"company_id"`
	Year            time.Time `json:"year"`
	Sales           float64   `json:"sales"`
	OperatingProfit float64   `json:"operating_profit"`
	NetProfit       float64   `json:"net_profit"`
	EPS             float64   `json:"eps"`
	EquityCapital   float64   `json:"equity_capital"`
	Reserves        float64   `json:"reserves"`
	Borrowings      float64   `json:"borrowings"`
	TotalAssets     float64   `json:"total_assets"`
}
