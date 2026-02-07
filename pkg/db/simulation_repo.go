package db

import (
	"context"
	"fmt"
	"time"

	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (s *Store) SaveSimulationState(ctx context.Context, state models.SimulationState) error {
	_, err := s.pool.Exec(
		ctx,
		`INSERT INTO simulation_state (id, tick, is_active, is_paused, start_time, last_update)
		 VALUES (1, $1, $2, $3, $4, $5)
		 ON CONFLICT (id) DO UPDATE SET
			tick = EXCLUDED.tick,
			is_active = EXCLUDED.is_active,
			is_paused = EXCLUDED.is_paused,
			start_time = EXCLUDED.start_time,
			last_update = EXCLUDED.last_update`,
		state.Tick,
		state.IsActive,
		state.IsPaused,
		state.StartTime,
		time.Now(),
	)
	if err != nil {
		logs.Log.Error("Failed to save simulation state", zap.Any("state", state), zap.Error(err))
		return fmt.Errorf("could not save simulation state: %w", err)
	}
	logs.Log.Debug("Simulation state saved", zap.Int("tick", state.Tick))
	return nil
}

func (s *Store) LoadSimulationState(ctx context.Context) (*models.SimulationState, error) {
	var state models.SimulationState
	err := s.pool.QueryRow(
		ctx,
		"SELECT id, tick, is_active, is_paused, start_time, last_update FROM simulation_state WHERE id = 1",
	).Scan(
		&state.ID,
		&state.Tick,
		&state.IsActive,
		&state.IsPaused,
		&state.StartTime,
		&state.LastUpdate,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			logs.Log.Info("No simulation state found in database")
			return nil, nil // No state persisted yet
		}
		logs.Log.Error("Failed to load simulation state", zap.Error(err))
		return nil, fmt.Errorf("could not load simulation state: %w", err)
	}
	logs.Log.Debug("Simulation state loaded", zap.Int("tick", state.Tick))
	return &state, nil
}
