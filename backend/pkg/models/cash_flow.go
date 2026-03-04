package models

import (
	"time"
)

type CashFlow struct {
	ID                        int       `json:"id"`
	CompanyID                 int       `json:"company_id"`
	Year                      time.Time `json:"year"`
	CashFromOperatingActivity float64   `json:"cash_from_operating_activity"`
	CashFromInvestingActivity float64   `json:"cash_from_investing_activity"`
	CashFromFinancingActivity float64   `json:"cash_from_financing_activity"`
	NetCashFlow               float64   `json:"net_cash_flow"`
}
