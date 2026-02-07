package db

import (
	"context"

	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"go.uber.org/zap"
)

func (s *Store) GetUser(ctx context.Context, username string, email string) (*models.User, error) {
	var user models.User

	err := s.pool.QueryRow(
		ctx,
		`SELECT id, username, email, cash_balance
		 FROM users
		 WHERE username = $1 OR email = $2`,
		username,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CashBalance,
	)

	if err != nil {
		logs.Log.Error("Failed to get user", zap.String("username", username), zap.String("email", email), zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserPortfolio(ctx context.Context, userId int) ([]models.PortfolioEntry, error) {
	rows, err := s.pool.Query(
		ctx,
		`SELECT user_id, company_id, quantity
		 FROM portfolio_entries
		 WHERE user_id = $1`,
		userId,
	)

	if err != nil {
		logs.Log.Error("Failed to query user portfolio", zap.Int("user_id", userId), zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	portfolio := make([]models.PortfolioEntry, 0)

	for rows.Next() {
		var p models.PortfolioEntry
		if err := rows.Scan(&p.UserID, &p.CompanyID, &p.Quantity); err != nil {
			logs.Log.Error("Failed to scan portfolio row", zap.Int("user_id", userId), zap.Error(err))
			return nil, err
		}
		portfolio = append(portfolio, p)
	}

	if err := rows.Err(); err != nil {
		logs.Log.Error("Rows error in GetUserPortfolio", zap.Int("user_id", userId), zap.Error(err))
		return nil, err
	}

	return portfolio, nil
}

func (s *Store) GetUserTrades(ctx context.Context, userId int) ([]models.Trade, error) {
	rows, err := s.pool.Query(
		ctx,
		`SELECT id, user_id, company_id, quantity, trade_type, date
		 FROM trades
		 WHERE user_id = $1`,
		userId,
	)

	if err != nil {
		logs.Log.Error("Failed to query user trades", zap.Int("user_id", userId), zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	transactions := make([]models.Trade, 0)

	for rows.Next() {
		var t models.Trade
		if err := rows.Scan(&t.ID, &t.UserID, &t.CompanyID, &t.Quantity, &t.TradeType, &t.Date); err != nil {
			logs.Log.Error("Failed to scan trade row", zap.Int("user_id", userId), zap.Error(err))
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		logs.Log.Error("Rows error in GetUserTrades", zap.Int("user_id", userId), zap.Error(err))
		return nil, err
	}

	return transactions, nil
}
