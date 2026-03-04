package db

import (
	"context"
	"fmt"
	"time"

	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Store struct {
	pool               *pgxpool.Pool
	CacheStocks        map[int][]models.StockPrice
	CachePF            map[int][]models.ProfitLoss
	CacheCF            map[int][]models.CashFlow
	CacheBS            map[int][]models.BalanceSheet
	CacheRatios        map[int][]models.Ratio
	CacheCurrentStocks []models.StockPrice
	LastCachedYear     int
}

func NewStore(databaseURL string) (*Store, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	return NewStoreWithConfig(config)
}

func NewStoreWithConfig(config *pgxpool.Config) (*Store, error) {
	var pool *pgxpool.Pool
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		logs.Log.Info("Attempting to connect to database", zap.Int("attempt", i+1))
		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			err = pool.Ping(context.Background())
			if err == nil {
				logs.Log.Info("Successfully connected to database")
				s := &Store{
					pool:        pool,
					CacheStocks: make(map[int][]models.StockPrice),
					CachePF:     make(map[int][]models.ProfitLoss),
					CacheCF:     make(map[int][]models.CashFlow),
					CacheBS:     make(map[int][]models.BalanceSheet),
					CacheRatios: make(map[int][]models.Ratio),
				}
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

	CREATE TABLE IF NOT EXISTS drawings (
		company_id INTEGER PRIMARY KEY,
		data JSONB NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
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

func (s *Store) UpdateStore(ctx context.Context, tick int) error {
	// 1. Get current stocks for this tick
	currentStocks, err := s.GetStocksAtTick(ctx, tick)
	if err != nil {
		return err
	}
	s.CacheCurrentStocks = currentStocks

	if len(currentStocks) == 0 {
		return nil
	}

	// 2. Update stock history cache efficiently
	for _, stock := range currentStocks {
		history, ok := s.CacheStocks[stock.CompanyID]
		// If cache is missing OR there's a jump in ticks (history length != current tick),
		// we fetch the full history till now.
		if !ok || len(history) != tick {
			fullHistory, err := s.GetStocksTillTickOfCompany(ctx, tick, stock.CompanyID)
			if err == nil {
				s.CacheStocks[stock.CompanyID] = fullHistory
			} else {
				logs.Log.Error("Failed to fetch initial/refreshed stock history", zap.Int("company_id", stock.CompanyID), zap.Int("tick", tick), zap.Error(err))
				s.CacheStocks[stock.CompanyID] = []models.StockPrice{stock}
			}
		} else {
			// Normal case: append the new stock price
			if len(history) == 0 || history[len(history)-1].Date.Before(stock.Date) {
				s.CacheStocks[stock.CompanyID] = append(history, stock)
			}
		}
	}

	// 3. Check if year changed to update fundamental data
	currentYear := currentStocks[0].Date.Year()
	if s.LastCachedYear == currentYear && len(s.CachePF) > 0 {
		return nil // No need to update fundamentals
	}

	logs.Log.Info("Year changed or cache empty, updating fundamentals", zap.Int("year", currentYear))

	// Update fundamentals for all companies.
	// Acquire a SINGLE connection for the entire loop so that all ~4×N queries
	// share one pool slot instead of each borrowing their own connection.
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection for fundamentals update: %w", err)
	}
	defer conn.Release()

	companies, err := s.GetAllCompanies(ctx)
	if err != nil {
		return err
	}

	date := currentStocks[0].Date
	for _, company := range companies {
		// We ignore errors here to ensure one faulty company doesn't stop the whole update
		bs, _ := s.getBalanceSheetsConn(ctx, conn, company.ID, date)
		cf, _ := s.getCashFlowsConn(ctx, conn, company.ID, date)
		pf, _ := s.getProfitLossConn(ctx, conn, company.ID, date)
		ratios, _ := s.getRatiosConn(ctx, conn, company.ID, date)

		s.CacheBS[company.ID] = bs
		s.CacheCF[company.ID] = cf
		s.CachePF[company.ID] = pf
		s.CacheRatios[company.ID] = ratios
	}

	s.LastCachedYear = currentYear
	return nil
}
