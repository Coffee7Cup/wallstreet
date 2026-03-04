package models

import (
	"time"
)

type BalanceSheet struct {
	ID               int       `json:"id"`
	CompanyID        int       `json:"company_id"`
	Year             time.Time `json:"year"`
	EquityCapital    float64   `json:"equity_capital"`
	Reserves         float64   `json:"reserves"`
	Borrowings       float64   `json:"borrowings"`
	OtherLiabilities float64   `json:"other_liabilities"`
	TotalLiabilities float64   `json:"total_liabilities"`
	FixedAssets      float64   `json:"fixed_assets"`
	CWIP             float64   `json:"cwip"`
	Investments      float64   `json:"investments"`
	OtherAssets      float64   `json:"other_assets"`
	TotalAssets      float64   `json:"total_assets"`
}
