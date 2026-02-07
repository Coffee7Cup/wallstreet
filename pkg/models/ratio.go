package models

import (
	"time"
)

type Ratio struct {
	ID             int       `json:"id"`
	CompanyID      int       `json:"company_id"`
	Year           time.Time `json:"year"`
	ROE            float64   `json:"roe"`
	DebtEquity     float64   `json:"debt_equity"`
	OPM            float64   `json:"opm"`
	IntrinsicValue float64   `json:"intrinsic_value"`
}
