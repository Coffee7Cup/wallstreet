package models

type PortfolioEntry struct {
	UserID        int     `json:"user_id"`
	CompanyID     int     `json:"company_id"`
	CompanyName   string  `json:"company_name"`
	CompanySymbol string  `json:"company_symbol"`
	Quantity      int     `json:"quantity"`
	CurrentPrice  float64 `json:"current_price"`
}
