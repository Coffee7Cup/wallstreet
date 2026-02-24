package models

import (
	"time"
)

type Ratio struct {
	ID                  int       `json:"id"`
	CompanyID           int       `json:"company_id"`
	Year                time.Time `json:"year"`
	OPM                 float64   `json:"opm"`
	DebtorDays          int       `json:"debtor_days"`
	InventoryDays       int       `json:"inventory_days"`
	DaysPayable         int       `json:"days_payable"`
	CashConversionCycle int       `json:"cash_conversion_cycle"`
	WorkingCapitalDays  int       `json:"working_capital_days"`
	ROCEPercent         float64   `json:"roce_percent"`
}
