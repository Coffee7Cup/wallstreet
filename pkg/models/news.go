package models

import (
	"time"
)

type News struct {
	ID           int       `json:"id"`
	ReleaseDate  time.Time `json:"release_date"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	NewsType     string    `json:"news_type"`
	ImpactFactor float64   `json:"impact_factor"`
	Tick         int       `json:"tick"`
}
