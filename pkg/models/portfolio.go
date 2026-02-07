package models

type PortfolioEntry struct {
	UserID    int `json:"user_id"`
	CompanyID int `json:"company_id"`
	Quantity  int `json:"quantity"`
}
