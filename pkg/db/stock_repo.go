package db

import (
	"context"
	"fmt"
	"time"

	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"go.uber.org/zap"
)

func (s *Store) GetRatios(ctx context.Context, companyID int) ([]models.Ratio, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, company_id, year, roe, debt_equity, opm, intrinsic_value FROM ratios WHERE company_id = $1 ORDER BY year DESC", companyID)
	if err != nil {
		logs.Log.Error("Failed to query ratios", zap.Int("company_id", companyID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var ratios []models.Ratio
	for rows.Next() {
		var r models.Ratio
		err := rows.Scan(&r.ID, &r.CompanyID, &r.Year, &r.ROE, &r.DebtEquity, &r.OPM, &r.IntrinsicValue)
		if err != nil {
			logs.Log.Error("Failed to scan ratio row", zap.Int("company_id", companyID), zap.Error(err))
			return nil, err
		}
		ratios = append(ratios, r)
	}
	return ratios, nil
}

func (s *Store) GetFundamentals(ctx context.Context, companyID int) ([]models.Fundamental, error) {

	rows, err := s.pool.Query(ctx, `
		SELECT id, company_id, year, sales, operating_profit, net_profit, eps, equity_capital, reserves, borrowings, total_assets 
		FROM fundamentals 
		WHERE company_id = $1 
		ORDER BY year DESC
	`, companyID)
	if err != nil {
		logs.Log.Error("Failed to query fundamentals", zap.Int("company_id", companyID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var fundamentals []models.Fundamental
	for rows.Next() {
		var f models.Fundamental
		err := rows.Scan(
			&f.ID, &f.CompanyID, &f.Year, &f.Sales, &f.OperatingProfit,
			&f.NetProfit, &f.EPS, &f.EquityCapital, &f.Reserves,
			&f.Borrowings, &f.TotalAssets,
		)
		if err != nil {
			logs.Log.Error("Failed to scan fundamental row", zap.Int("company_id", companyID), zap.Error(err))
			return nil, err
		}
		fundamentals = append(fundamentals, f)
	}
	return fundamentals, nil
}

func (s *Store) GetStocksAtDate(ctx context.Context, date time.Time) ([]models.StockPrice, error) {
	rows, err := s.pool.Query(ctx, "SELECT * FROM stock_prices WHERE date = $1", date)
	if err != nil {
		logs.Log.Error("Failed to query stocks at date", zap.Time("date", date), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var stocks []models.StockPrice
	for rows.Next() {
		var stock models.StockPrice
		if err := rows.Scan(&stock.ID, &stock.CompanyID, &stock.Date, &stock.OpenPrice, &stock.ClosePrice, &stock.HighPrice, &stock.LowPrice, &stock.NoOfShares, &stock.NoOfTrades, &stock.TotalTurnover); err != nil {
			logs.Log.Error("Failed to scan stock row", zap.Time("date", date), zap.Error(err))
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func (s *Store) GetStocksAtDateOfCompany(ctx context.Context, date time.Time, companyID int) (models.StockPrice, error) {
	var stock models.StockPrice
	err := s.pool.QueryRow(ctx, "SELECT * FROM stock_prices WHERE date = $1 AND company_id = $2", date, companyID).Scan(&stock.ID, &stock.CompanyID, &stock.Date, &stock.OpenPrice, &stock.ClosePrice, &stock.HighPrice, &stock.LowPrice, &stock.NoOfShares, &stock.NoOfTrades, &stock.TotalTurnover)
	if err != nil {
		logs.Log.Error("Failed to get stock at date for company", zap.Time("date", date), zap.Int("company_id", companyID), zap.Error(err))
		return stock, err
	}
	return stock, nil
}

func (s *Store) GetStocksTillDateOfCompany(ctx context.Context, date time.Time, companyID int) ([]models.StockPrice, error) {
	rows, err := s.pool.Query(ctx, "SELECT * FROM stock_prices WHERE date <= $1 AND company_id = $2", date, companyID)
	if err != nil {
		logs.Log.Error("Failed to query stocks till date for company", zap.Time("date", date), zap.Int("company_id", companyID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var stocks []models.StockPrice
	for rows.Next() {
		var stock models.StockPrice
		if err := rows.Scan(&stock.ID, &stock.CompanyID, &stock.Date, &stock.OpenPrice, &stock.ClosePrice, &stock.HighPrice, &stock.LowPrice, &stock.NoOfShares, &stock.NoOfTrades, &stock.TotalTurnover); err != nil {
			logs.Log.Error("Failed to scan stock till date row", zap.Time("date", date), zap.Int("company_id", companyID), zap.Error(err))
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	if rows.Err() != nil {
		logs.Log.Error("Rows error in GetStocksTillDateOfCompany", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	return stocks, nil
}

func (s *Store) GetStocksTillDate(ctx context.Context, date time.Time) ([]models.StockPrice, error) {
	rows, err := s.pool.Query(ctx, "SELECT * FROM stock_prices WHERE date <= $1", date)
	if err != nil {
		logs.Log.Error("Failed to query stocks till date", zap.Time("date", date), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var stocks []models.StockPrice
	for rows.Next() {
		var stock models.StockPrice
		if err := rows.Scan(&stock.ID, &stock.CompanyID, &stock.Date, &stock.OpenPrice, &stock.ClosePrice, &stock.HighPrice, &stock.LowPrice, &stock.NoOfShares, &stock.NoOfTrades, &stock.TotalTurnover); err != nil {
			logs.Log.Error("Failed to scan stock till date row", zap.Time("date", date), zap.Error(err))
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func (s *Store) GetStocksAtTick(ctx context.Context, tick int) ([]models.StockPrice, error) {
	var stocks []models.StockPrice
	var err error

	for i := 0; i < 3; i++ {
		rows, errQuery := s.pool.Query(ctx, "SELECT id, company_id, date, open_price, close_price, high_price, low_price, no_of_shares, no_of_trades, total_turnover FROM stock_prices_with_ticks WHERE tick_idx = $1", tick)
		if errQuery != nil {
			err = errQuery
			logs.Log.Warn("Failed to query stocks at tick, retrying...", zap.Int("tick", tick), zap.Int("attempt", i+1), zap.Error(err))
			time.Sleep(100 * time.Millisecond)
			continue
		}
		defer rows.Close()

		stocks = nil // Clear if retry
		for rows.Next() {
			var stock models.StockPrice
			if errScan := rows.Scan(&stock.ID, &stock.CompanyID, &stock.Date, &stock.OpenPrice, &stock.ClosePrice, &stock.HighPrice, &stock.LowPrice, &stock.NoOfShares, &stock.NoOfTrades, &stock.TotalTurnover); errScan != nil {
				err = errScan
				logs.Log.Error("Failed to scan stock row at tick", zap.Int("tick", tick), zap.Error(err))
				return nil, err
			}
			stocks = append(stocks, stock)
		}
		return stocks, nil
	}

	return nil, fmt.Errorf("failed to get stocks at tick after 3 attempts: %w", err)
}

func (s *Store) GetStockAtTickOfCompany(ctx context.Context, tick int, companyID int) (models.StockPrice, error) {
	var stock models.StockPrice
	err := s.pool.QueryRow(ctx, "SELECT id, company_id, date, open_price, close_price, high_price, low_price, no_of_shares, no_of_trades, total_turnover FROM stock_prices_with_ticks WHERE tick_idx = $1 AND company_id = $2", tick, companyID).Scan(&stock.ID, &stock.CompanyID, &stock.Date, &stock.OpenPrice, &stock.ClosePrice, &stock.HighPrice, &stock.LowPrice, &stock.NoOfShares, &stock.NoOfTrades, &stock.TotalTurnover)
	if err != nil {
		logs.Log.Error("Failed to get stock at tick for company", zap.Int("tick", tick), zap.Int("company_id", companyID), zap.Error(err))
		return stock, err
	}
	return stock, nil
}

func (s *Store) GetStocksTillTickOfCompany(ctx context.Context, tick int, companyID int) ([]models.StockPrice, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, company_id, date, open_price, close_price, high_price, low_price, no_of_shares, no_of_trades, total_turnover FROM stock_prices_with_ticks WHERE tick_idx <= $1 AND company_id = $2 ORDER BY tick_idx ASC", tick, companyID)
	if err != nil {
		logs.Log.Error("Failed to query stocks till tick for company", zap.Int("tick", tick), zap.Int("company_id", companyID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var stocks []models.StockPrice
	for rows.Next() {
		var stock models.StockPrice
		if err := rows.Scan(&stock.ID, &stock.CompanyID, &stock.Date, &stock.OpenPrice, &stock.ClosePrice, &stock.HighPrice, &stock.LowPrice, &stock.NoOfShares, &stock.NoOfTrades, &stock.TotalTurnover); err != nil {
			logs.Log.Error("Failed to scan stock till tick row", zap.Int("tick", tick), zap.Int("company_id", companyID), zap.Error(err))
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	if rows.Err() != nil {
		logs.Log.Error("Rows error in GetStocksTillTickOfCompany", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	return stocks, nil
}

func (s *Store) GetStocksTillTick(ctx context.Context, tick int) ([]models.StockPrice, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, company_id, date, open_price, close_price, high_price, low_price, no_of_shares, no_of_trades, total_turnover FROM stock_prices_with_ticks WHERE tick_idx <= $1 ORDER BY tick_idx ASC", tick)
	if err != nil {
		logs.Log.Error("Failed to query stocks till tick for company", zap.Int("tick", tick), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var stocks []models.StockPrice
	for rows.Next() {
		var stock models.StockPrice
		if err := rows.Scan(&stock.ID, &stock.CompanyID, &stock.Date, &stock.OpenPrice, &stock.ClosePrice, &stock.HighPrice, &stock.LowPrice, &stock.NoOfShares, &stock.NoOfTrades, &stock.TotalTurnover); err != nil {
			logs.Log.Error("Failed to scan stock till tick row", zap.Int("tick", tick), zap.Error(err))
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	if rows.Err() != nil {
		logs.Log.Error("Rows error in GetStocksTillTickOfCompany", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	return stocks, nil
}

func (s *Store) GetAllCompanies(ctx context.Context) ([]models.Company, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, symbol, name, sector, total_shares FROM companies")
	if err != nil {
		logs.Log.Error("Failed to query companies", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var companies []models.Company
	for rows.Next() {
		var c models.Company
		if err := rows.Scan(&c.ID, &c.Symbol, &c.Name, &c.Sector, &c.TotalShares); err != nil {
			logs.Log.Error("Failed to scan company row", zap.Error(err))
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}

func (s *Store) GetCompanyByID(ctx context.Context, id int) (models.Company, error) {
	var c models.Company
	err := s.pool.QueryRow(ctx, "SELECT id, symbol, name, sector, total_shares FROM companies WHERE id = $1", id).Scan(&c.ID, &c.Symbol, &c.Name, &c.Sector, &c.TotalShares)
	if err != nil {
		logs.Log.Error("Failed to get company by ID", zap.Int("id", id), zap.Error(err))
		return c, err
	}
	return c, nil
}
