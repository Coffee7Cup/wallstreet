package db

import (
	"context"
	"fmt"
	"time"

	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"go.uber.org/zap"
)

func (s *Store) GetNewsByDay(ctx context.Context, date time.Time) ([]models.News, error) {
	rows, err := s.pool.Query(ctx, "SELECT n.id, n.release_date, n.title, n.content, n.company_id, c.name, c.symbol FROM news n JOIN companies c ON n.company_id = c.id WHERE n.release_date = $1", date)
	if err != nil {
		logs.Log.Error("Failed to query news", zap.Time("date", date), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(
			&n.ID,
			&n.ReleaseDate,
			&n.Title,
			&n.Content,
			&n.CompanyId,
			&n.CompanyName,
			&n.CompanySymbol,
		)
		if err != nil {
			logs.Log.Error("Failed to scan news row", zap.Time("date", date), zap.Error(err))
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}

func (s *Store) GetNewsByDateAndCompanyId(ctx context.Context, date time.Time, companyId int) ([]models.News, error) {
	rows, err := s.pool.Query(ctx, "SELECT n.id, n.release_date, n.title, n.content, n.company_id, c.name, c.symbol FROM news n JOIN companies c ON n.company_id = c.id WHERE n.release_date = $1 AND n.company_id = $2", date, companyId)
	if err != nil {
		logs.Log.Error("Failed to query news", zap.Time("date", date), zap.Int("company_id", companyId), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(
			&n.ID,
			&n.ReleaseDate,
			&n.Title,
			&n.Content,
			&n.CompanyId,
			&n.CompanyName,
			&n.CompanySymbol,
		)
		if err != nil {
			logs.Log.Error("Failed to scan news row", zap.Time("date", date), zap.Int("company_id", companyId), zap.Error(err))
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}

func (s *Store) GetNewsTillDateByCompanyId(ctx context.Context, date time.Time, companyId int) ([]models.News, error) {
	rows, err := s.pool.Query(ctx, "SELECT n.id, n.release_date, n.title, n.content, n.company_id, c.name, c.symbol FROM news n JOIN companies c ON n.company_id = c.id WHERE n.release_date <= $1 AND n.company_id = $2", date, companyId)
	if err != nil {
		logs.Log.Error("Failed to query news", zap.Time("date", date), zap.Int("company_id", companyId), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(
			&n.ID,
			&n.ReleaseDate,
			&n.Title,
			&n.Content,
			&n.CompanyId,
			&n.CompanyName,
			&n.CompanySymbol,
		)
		if err != nil {
			logs.Log.Error("Failed to scan news row", zap.Time("date", date), zap.Int("company_id", companyId), zap.Error(err))
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}

func (s *Store) GetNewsSearchLimitOffsetByCompany(ctx context.Context, companyID int, date time.Time, limit int, offset int, query string) ([]models.News, error) {
	sql := `
		SELECT n.id, n.release_date, n.title, n.content, n.company_id, c.name, c.symbol FROM news n
        JOIN companies c ON n.company_id = c.id
		WHERE n.company_id = $1 AND n.release_date <= $2
	`
	args := []interface{}{companyID, date}
	argIdx := 3

	if query != "" {
		sql += fmt.Sprintf(" AND (n.title ILIKE $%d OR n.content ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+query+"%")
		argIdx++
	}

	sql += fmt.Sprintf(" ORDER BY n.release_date DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		logs.Log.Error("Failed to query news with search/limit", zap.Int("company_id", companyID), zap.Time("date", date), zap.String("query", query), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(&n.ID, &n.ReleaseDate, &n.Title, &n.Content, &n.CompanyId, &n.CompanyName, &n.CompanySymbol)
		if err != nil {
			logs.Log.Error("Failed to scan news row with search/limit", zap.Error(err))
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}

func (s *Store) GetNewsAtTickForAll(ctx context.Context, tick int) ([]models.News, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT n.id, n.release_date, n.title, n.content, n.company_id, c.name, c.symbol
		FROM news n
        JOIN companies c ON n.company_id = c.id
		JOIN stock_prices_with_ticks spt ON n.company_id = spt.company_id AND n.release_date = spt.date
		WHERE spt.tick_idx = $1
	`, tick)
	if err != nil {
		logs.Log.Error("Failed to query news at tick for all", zap.Int("tick", tick), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(
			&n.ID,
			&n.ReleaseDate,
			&n.Title,
			&n.Content,
			&n.CompanyId,
			&n.CompanyName,
			&n.CompanySymbol,
		)
		if err != nil {
			logs.Log.Error("Failed to scan news row at tick", zap.Error(err))
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}
func (s *Store) GetNewsTillTickLimitSearchForAll(ctx context.Context, tick int, limit int, offset int, query string) ([]models.News, error) {
	sql := `
		SELECT n.id, n.release_date, n.title, n.content, n.company_id, c.name, c.symbol
		FROM news n
		JOIN companies c ON n.company_id = c.id
		JOIN stock_prices_with_ticks spt ON n.company_id = spt.company_id AND n.release_date = spt.date
		WHERE spt.tick_idx <= $1
	`
	args := []interface{}{tick}
	argIdx := 2

	if query != "" {
		sql += fmt.Sprintf(" AND (n.title ILIKE $%d OR n.content ILIKE $%d OR c.name ILIKE $%d OR c.symbol ILIKE $%d)", argIdx, argIdx, argIdx, argIdx)
		args = append(args, "%"+query+"%")
		argIdx++
	}

	sql += fmt.Sprintf(" ORDER BY n.release_date DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		logs.Log.Error("Failed to query news till tick for all", zap.Int("tick", tick), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(&n.ID, &n.ReleaseDate, &n.Title, &n.Content, &n.CompanyId, &n.CompanyName, &n.CompanySymbol)
		if err != nil {
			logs.Log.Error("Failed to scan news row till tick", zap.Error(err))
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}
func (s *Store) GetNewsBySectorTillTick(ctx context.Context, sector string, tick int, limit int, offset int) ([]models.News, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT n.id, n.release_date, n.title, n.content, n.company_id, c.name, c.symbol
		FROM news n
		JOIN companies c ON n.company_id = c.id
		JOIN stock_prices_with_ticks spt ON n.company_id = spt.company_id AND n.release_date = spt.date
		WHERE spt.tick_idx <= $1 AND c.sector = $2
		ORDER BY n.release_date DESC
		LIMIT $3 OFFSET $4
	`, tick, sector, limit, offset)
	if err != nil {
		logs.Log.Error("Failed to query news by sector till tick", zap.String("sector", sector), zap.Int("tick", tick), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var newsList []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(&n.ID, &n.ReleaseDate, &n.Title, &n.Content, &n.CompanyId, &n.CompanyName, &n.CompanySymbol)
		if err != nil {
			logs.Log.Error("Failed to scan news row by sector", zap.Error(err))
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}
