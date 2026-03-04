package models

import (
	"time"
)

type News struct {
	ID            int       `json:"id"`
	ReleaseDate   time.Time `json:"release_date"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CompanyId     int       `json:"company_id"`
	CompanyName   string    `json:"company_name"`
	CompanySymbol string    `json:"company_symbol"`
}
