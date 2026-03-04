package db

import (
	"context"
	"fmt"

	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (s *Store) SellStocks(ctx context.Context, trade models.Trade) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		logs.Log.Error("Failed to start transaction for selling stocks", zap.Any("trade", trade), zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)
	logs.Log.Debug("Started transaction for selling stocks", zap.Int("user_id", trade.UserID), zap.Int("company_id", trade.CompanyID))

	// 1. Fetch Authoritative Price
	var actualPrice float64
	err = tx.QueryRow(ctx, "SELECT close_price FROM stock_prices WHERE company_id = $1 AND date = $2", trade.CompanyID, trade.Date).Scan(&actualPrice)
	if err != nil {
		logs.Log.Error("Failed to fetch authoritative price for selling", zap.Any("trade", trade), zap.Error(err))
		return err
	}

	// 2. Lock User Cash First (same order as BuyStocks: users → portfolio_entries)
	// This prevents deadlocks when BUY and SELL run concurrently for the same user.
	var cash float64
	err = tx.QueryRow(
		ctx,
		"SELECT cash_balance FROM users WHERE id = $1 FOR UPDATE",
		trade.UserID,
	).Scan(&cash)
	if err != nil {
		logs.Log.Error("Failed to lock user row for selling", zap.Any("trade", trade), zap.Error(err))
		return err
	}

	// 3. Check and Lock Portfolio Entry
	var ownedQty int
	err = tx.QueryRow(
		ctx,
		"SELECT quantity FROM portfolio_entries WHERE user_id = $1 AND company_id = $2 FOR UPDATE",
		trade.UserID, trade.CompanyID,
	).Scan(&ownedQty)

	if err != nil {
		logs.Log.Error("Failed to check portfolio for selling", zap.Any("trade", trade), zap.Error(err))
		return err
	}

	if trade.Quantity > ownedQty {
		err := fmt.Errorf("not enough shares: you own %d but tried to sell %d", ownedQty, trade.Quantity)
		logs.Log.Warn("Insufficient shares for sale", zap.Int("user_id", trade.UserID), zap.Int("company_id", trade.CompanyID), zap.Int("owned", ownedQty), zap.Int("requested", trade.Quantity))
		return err
	}

	// 4. Update Portfolio
	_, err = tx.Exec(
		ctx,
		"UPDATE portfolio_entries SET quantity = quantity - $1 WHERE user_id = $2 AND company_id = $3",
		trade.Quantity, trade.UserID, trade.CompanyID,
	)
	if err != nil {
		logs.Log.Error("Failed to update portfolio for selling", zap.Any("trade", trade), zap.Error(err))
		return err
	}

	// 5. Update Cash Balance
	revenue := actualPrice * float64(trade.Quantity)
	_, err = tx.Exec(
		ctx,
		"UPDATE users SET cash_balance = cash_balance + $1 WHERE id = $2",
		revenue, trade.UserID,
	)
	if err != nil {
		logs.Log.Error("Failed to update cash balance for selling", zap.Any("trade", trade), zap.Error(err))
		return err
	}

	// 5. Log Trade Record
	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO trades (user_id, company_id, trade_type, quantity, price, date, timestamp)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		`,
		trade.UserID, trade.CompanyID, trade.TradeType,
		trade.Quantity, actualPrice, trade.Date, trade.Timestamp,
	)
	if err != nil {
		logs.Log.Error("Failed to log trade record for selling", zap.Any("trade", trade), zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logs.Log.Error("Failed to commit transaction for selling stocks", zap.Any("trade", trade), zap.Error(err))
		return err
	}
	logs.Log.Info("Successfully sold stocks", zap.Int("user_id", trade.UserID), zap.Int("company_id", trade.CompanyID), zap.Int("quantity", trade.Quantity))
	return nil
}

// i guess ii must take date from the engine state ***
func (s *Store) BuyStocks(ctx context.Context, trade models.Trade) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		logs.Log.Error("Failed to start transaction for buying stocks", zap.Any("trade", trade), zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)
	logs.Log.Debug("Started transaction for buying stocks", zap.Int("user_id", trade.UserID), zap.Int("company_id", trade.CompanyID))

	// 1. Fetch Authoritative Price
	var actualPrice float64
	err = tx.QueryRow(ctx, "SELECT close_price FROM stock_prices WHERE company_id = $1 AND date = $2", trade.CompanyID, trade.Date).Scan(&actualPrice)
	if err != nil {
		if err == pgx.ErrNoRows {
			logs.Log.Warn("Stock price not found for buying", zap.Int("company_id", trade.CompanyID), zap.Time("date", trade.Date))
			return fmt.Errorf("stock price not found for company %d on date %s", trade.CompanyID, trade.Date.Format("2006-01-02"))
		}
		logs.Log.Error("Failed to fetch authoritative price for buying", zap.Any("trade", trade), zap.Error(err))
		return fmt.Errorf("could not fetch authoritative price: %w", err)
	}

	// 2. Check and Lock User Cash
	var cash float64
	err = tx.QueryRow(
		ctx,
		"SELECT cash_balance FROM users WHERE id = $1 FOR UPDATE",
		trade.UserID,
	).Scan(&cash)
	if err != nil {
		logs.Log.Error("Failed to check user cash for buying", zap.Any("trade", trade), zap.Error(err))
		return err
	}

	cost := actualPrice * float64(trade.Quantity)
	if cash < cost {
		err := fmt.Errorf("not enough cash: cost %.2f, balance %.2f", cost, cash)
		logs.Log.Warn("Insufficient cash for buying", zap.Int("user_id", trade.UserID), zap.Float64("balance", cash), zap.Float64("cost", cost))
		return err
	}

	// 3. Update Cash
	_, err = tx.Exec(ctx, "UPDATE users SET cash_balance = cash_balance - $1 WHERE id = $2", cost, trade.UserID)
	if err != nil {
		logs.Log.Error("Failed to update cash for buying", zap.Any("trade", trade), zap.Error(err))
		return err
	}

	// 4. Update Portfolio
	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO portfolio_entries (user_id, company_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, company_id)
		DO UPDATE SET quantity = portfolio_entries.quantity + $3
		`,
		trade.UserID, trade.CompanyID, trade.Quantity,
	)
	if err != nil {
		logs.Log.Error("Failed to update portfolio for buying", zap.Any("trade", trade), zap.Error(err))
		return err
	}

	// 5. Log Trade Record
	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO trades (user_id, company_id, trade_type, quantity, price, date, timestamp)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		`,
		trade.UserID, trade.CompanyID, trade.TradeType,
		trade.Quantity, actualPrice, trade.Date, trade.Timestamp,
	)
	if err != nil {
		logs.Log.Error("Failed to log trade record for buying", zap.Any("trade", trade), zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logs.Log.Error("Failed to commit transaction for buying stocks", zap.Any("trade", trade), zap.Error(err))
		return err
	}
	logs.Log.Info("Successfully bought stocks", zap.Int("user_id", trade.UserID), zap.Int("company_id", trade.CompanyID), zap.Int("quantity", trade.Quantity))
	return nil
}

func (s *Store) GetUserTradesLimit(ctx context.Context, userID int, limit int) ([]models.Trade, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT t.id, t.user_id, t.company_id, c.name, c.symbol, t.trade_type, t.quantity, t.price, t.date, t.timestamp 
		FROM trades t
		JOIN companies c ON t.company_id = c.id
		WHERE t.user_id = $1 
		ORDER BY t.timestamp DESC 
		LIMIT $2
	`, userID, limit)
	if err != nil {
		logs.Log.Error("Failed to fetch user trades", zap.Int("user_id", userID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var trades []models.Trade
	for rows.Next() {
		var trade models.Trade
		if err := rows.Scan(&trade.ID, &trade.UserID, &trade.CompanyID, &trade.CompanyName, &trade.CompanySymbol, &trade.TradeType, &trade.Quantity, &trade.Price, &trade.Date, &trade.Timestamp); err != nil {
			logs.Log.Error("Failed to scan trade", zap.Error(err))
			return nil, err
		}
		trades = append(trades, trade)
	}

	return trades, nil
}

func (s *Store) GetUserTradesSearchLimit(ctx context.Context, userID int, limit int, offset int, query string) ([]models.Trade, error) {
	sql := `
		SELECT t.id, t.user_id, t.company_id, c.name, c.symbol, t.trade_type, t.quantity, t.price, t.date, t.timestamp 
		FROM trades t
		JOIN companies c ON t.company_id = c.id
		WHERE t.user_id = $1
	`
	args := []interface{}{userID}
	argIdx := 2

	if query != "" {
		sql += fmt.Sprintf(" AND (c.symbol ILIKE $%d OR c.name ILIKE $%d OR t.trade_type ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+query+"%")
		argIdx++
	}

	sql += fmt.Sprintf(" ORDER BY t.timestamp DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		logs.Log.Error("Failed to fetch user trades with search/limit", zap.Int("user_id", userID), zap.String("query", query), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var trades []models.Trade
	for rows.Next() {
		var trade models.Trade
		if err := rows.Scan(&trade.ID, &trade.UserID, &trade.CompanyID, &trade.CompanyName, &trade.CompanySymbol, &trade.TradeType, &trade.Quantity, &trade.Price, &trade.Date, &trade.Timestamp); err != nil {
			logs.Log.Error("Failed to scan trade search/limit", zap.Error(err))
			return nil, err
		}
		trades = append(trades, trade)
	}

	return trades, nil
}

func (s *Store) GetAllSectors(ctx context.Context) ([]string, error) {
	sectors, err := s.pool.Query(ctx, "SELECT DISTINCT sector FROM companies")
	if err != nil {
		logs.Log.Error("Failed to fetch all sectors", zap.Error(err))
		return nil, err
	}
	defer sectors.Close()

	var sectorsList []string
	for sectors.Next() {
		var sector string
		if err := sectors.Scan(&sector); err != nil {
			logs.Log.Error("Failed to scan sector", zap.Error(err))
			return nil, err
		}
		sectorsList = append(sectorsList, sector)
	}

	return sectorsList, nil
}
