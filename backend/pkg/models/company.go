package models

type Company struct {
	ID          int    `json:"id"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Sector      string `json:"sector"`
	TotalShares int64  `json:"total_shares"`
}
