package models

import (
	"time"
)

type News struct {
	ID           int       `json:"id"`
	ReleaseDate  time.Time `json:"release_date"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Tick         int       `json:"tick"`
	CompanyId	 int 	   `json:"company_id"`
}
