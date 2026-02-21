package models

type User struct {
	ID          int     `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	CashBalance float64 `json:"cash_balance"`
}

type TopTrader struct {
	Rank        int     `json:"rank"`
	Username    string  `json:"username"`
	CashBalance float64 `json:"cash_balance"`
	StockValue  float64 `json:"stock_value"`
	TotalValue  float64 `json:"total_value"`
}
