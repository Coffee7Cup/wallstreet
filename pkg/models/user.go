package models

type User struct {
	ID          int     `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	CashBalance float64 `json:"cash_balance"`
}
