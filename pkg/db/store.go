package db

import (
	"context"
	"fmt"
	"time"

	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(databaseURL string) (*Store, error) {
	var pool *pgxpool.Pool
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		logs.Log.Info("Attempting to connect to database", zap.Int("attempt", i+1))
		pool, err = pgxpool.New(context.Background(), databaseURL)
		if err == nil {
			err = pool.Ping(context.Background())
			if err == nil {
				logs.Log.Info("Successfully connected to database")
				s := &Store{pool: pool}
				if err := s.RunInitSQL(context.Background()); err != nil {
					logs.Log.Error("Failed to run initialization SQL", zap.Error(err))
					return nil, err
				}
				return s, nil
			}
		}

		logs.Log.Warn("Failed to connect to database, retrying...", zap.Int("attempt", i+1), zap.Error(err))
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("unable to connect to database after %d attempts: %w", maxRetries, err)
}

func (s *Store) RunInitSQL(ctx context.Context) error {
	logs.Log.Info("Ensuring database views exist...")

	initSQL := `
	CREATE OR REPLACE VIEW stock_prices_with_ticks AS
	SELECT *, ROW_NUMBER() OVER (PARTITION BY company_id ORDER BY date) - 1 as tick_idx
	FROM stock_prices;
	`

	_, err := s.pool.Exec(ctx, initSQL)
	if err != nil {
		return fmt.Errorf("could not run init sql: %w", err)
	}

	logs.Log.Info("Database views ensured")
	return nil
}

func (s *Store) Close() {
	s.pool.Close()
}
