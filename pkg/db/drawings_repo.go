package db

import (
	"context"
	"encoding/json"
	"fmt"
)

func (s *Store) SaveDrawings(ctx context.Context, companyID int, drawings json.RawMessage) error {
	query := `
		INSERT INTO drawings (company_id, data, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (company_id)
		DO UPDATE SET data = EXCLUDED.data, updated_at = CURRENT_TIMESTAMP;
	`
	_, err := s.pool.Exec(ctx, query, companyID, drawings)
	if err != nil {
		return fmt.Errorf("failed to save drawings: %w", err)
	}
	return nil
}

func (s *Store) GetDrawings(ctx context.Context, companyID int) (json.RawMessage, error) {
	query := `SELECT data FROM drawings WHERE company_id = $1`
	var data json.RawMessage
	err := s.pool.QueryRow(ctx, query, companyID).Scan(&data)
	if err != nil {
		return nil, nil // Return nil if no drawings found, not an error
	}
	return data, nil
}
