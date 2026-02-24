package models

import (
	"time"
)

type ProfitLoss struct {
	ID              int       `json:"id"`
	CompanyID       int       `json:"company_id"`
	Year            time.Time `json:"year"`
	Sales           float64   `json:"sales"`
	Expenses        float64   `json:"expenses"`
	OperatingProfit float64   `json:"operating_profit"`
	OPMPercent      float64   `json:"opm_percent"`
	OtherIncome     float64   `json:"other_income"`
	Interest        float64   `json:"interest"`
	Depreciation    float64   `json:"depreciation"`
	ProfitBeforeTax float64   `json:"profit_before_tax"`
	TaxPercent      float64   `json:"tax_percent"`
	NetProfit       float64   `json:"net_profit"`
	EPS             float64   `json:"eps"`
	DividendPayout  float64   `json:"dividend_payout"`
}
