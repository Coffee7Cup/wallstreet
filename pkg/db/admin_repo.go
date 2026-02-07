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
