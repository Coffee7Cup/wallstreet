package models

import (
	"time"
)

type SimulationState struct {
	ID         int       `json:"id"`
	Tick       int       `json:"tick"`
	IsActive   bool      `json:"is_active"`
	IsPaused   bool      `json:"is_paused"`
	StartTime  time.Time `json:"start_time"`
	LastUpdate time.Time `json:"last_update"`
}

type SimulationBroadcast struct {
	Tick     int          `json:"tick"`
	IsActive bool         `json:"is_active"`
	IsPaused bool         `json:"is_paused"`
}
