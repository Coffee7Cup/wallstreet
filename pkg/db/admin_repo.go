package db

import (
	"context"

	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"go.uber.org/zap"
)

func (s *Store) GetAdmin(ctx context.Context, username string, email string) (*models.Admin, error) {
	var admin models.Admin
	err := s.pool.QueryRow(ctx,
		"SELECT id, username, email FROM admins WHERE username = $1 OR email = $2",
		username,
		email,
	).Scan(
		&admin.ID,
		&admin.Username,
		&admin.Email,
	)
	if err != nil {
		logs.Log.Error("Failed to get admin", zap.String("username", username), zap.String("email", email), zap.Error(err))
		return nil, err
	}
	return &admin, nil
}

func (s *Store) GetTradesByUserID(ctx context.Context, user_id string) ([]models.Trade, error) {
	var trades []models.Trade
	rows, err := s.pool.Query(ctx,
		"SELECT * FROM trades WHERE user_id = $1",
		user_id,
	)
	if err != nil {
		logs.Log.Error("Failed to get trades", zap.String("user_id", user_id), zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var trade models.Trade
		err := rows.Scan(
			&trade.ID,
			&trade.UserID,
			&trade.CompanyID,
			&trade.TradeType,
			&trade.Quantity,
			&trade.Date,
			&trade.Timestamp,
		)
		if err != nil {
			logs.Log.Error("Failed to scan trade", zap.Error(err))
			return nil, err
		}
		trades = append(trades, trade)
	}

	if rows.Err() != nil {
		logs.Log.Error("Failed to get trades", zap.String("user_id", user_id), zap.Error(rows.Err()))
		return nil, rows.Err()
	}
	return trades, nil
}

func (s *Store) GetTradesBySymbol(ctx context.Context, user_id string, symbol string) ([]models.Trade, error) {
	var trades []models.Trade
	rows, err := s.pool.Query(ctx,
		"SELECT * FROM trades WHERE user_id = $1 AND symbol = $2",
		user_id,
		symbol,
	)
	if err != nil {
		logs.Log.Error("Failed to get trades", zap.String("user_id", user_id), zap.String("symbol", symbol), zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var trade models.Trade
		err := rows.Scan(
			&trade.ID,
			&trade.UserID,
			&trade.CompanyID,
			&trade.TradeType,
			&trade.Quantity,
			&trade.Date,
			&trade.Timestamp,
		)
		if err != nil {
			logs.Log.Error("Failed to scan trade", zap.Error(err))
			return nil, err
		}
		trades = append(trades, trade)
	}

	if rows.Err() != nil {
		logs.Log.Error("Failed to get trades", zap.String("user_id", user_id), zap.String("symbol", symbol), zap.Error(rows.Err()))
		return nil, rows.Err()
	}
	return trades, nil
}

func (s *Store) GetTradesBySymbolAndLimit(ctx context.Context, user_id string, symbol string, limit int) ([]models.Trade, error) {
	var trades []models.Trade
	rows, err := s.pool.Query(ctx,
		"SELECT * FROM trades WHERE user_id = $1 AND symbol = $2 LIMIT $3",
		user_id,
		symbol,
		limit,
	)
	if err != nil {
		logs.Log.Error("Failed to get trades", zap.String("user_id", user_id), zap.String("symbol", symbol), zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var trade models.Trade
		err := rows.Scan(
			&trade.ID,
			&trade.UserID,
			&trade.CompanyID,
			&trade.TradeType,
			&trade.Quantity,
			&trade.Date,
			&trade.Timestamp,
		)
		if err != nil {
			logs.Log.Error("Failed to scan trade", zap.Error(err))
			return nil, err
		}
		trades = append(trades, trade)
	}

	if rows.Err() != nil {
		logs.Log.Error("Failed to get trades", zap.String("user_id", user_id), zap.String("symbol", symbol), zap.Error(rows.Err()))
		return nil, rows.Err()
	}
	return trades, nil
}

func (s *Store) ResetSimulation(ctx context.Context) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		logs.Log.Error("Failed to start transaction for resetting simulation", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM portfolio_entries")
	if err != nil {
		logs.Log.Error("Failed to delete portfolio_entries", zap.Error(err))
		return err
	}
	_, err = tx.Exec(ctx, "DELETE FROM trades")
	if err != nil {
		logs.Log.Error("Failed to delete trades", zap.Error(err))
		return err
	}
	_, err = tx.Exec(ctx, "DELETE FROM simulation_state")
	if err != nil {
		logs.Log.Error("Failed to delete simulation_state", zap.Error(err))
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE users SET cash_balance = 100000")
	if err != nil {
		logs.Log.Error("Failed to update users", zap.Error(err))
		return err
	}
	return tx.Commit(ctx)
}
