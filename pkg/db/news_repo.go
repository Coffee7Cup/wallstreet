package db

import (
	"context"

	"time"

	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"go.uber.org/zap"
)

func (s *Store) GetNewsByDay(ctx context.Context, date time.Time) ([]models.News, error) {
	rows, err := s.pool.Query(ctx, "SELECT * FROM news WHERE release_date = $1", date)
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
			&n.NewsType,
			&n.ImpactFactor,
		)
		if err != nil {
			logs.Log.Error("Failed to scan news row", zap.Time("date", date), zap.Error(err))
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}
func (s *Store) GetNewsAtTick(ctx context.Context, tick int) ([]models.News, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, release_date, title, content, news_type, impact_factor, tick_idx FROM news_with_ticks WHERE tick_idx = $1", tick)
	if err != nil {
		logs.Log.Error("Failed to query news at tick", zap.Int("tick", tick), zap.Error(err))
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
			&n.NewsType,
			&n.ImpactFactor,
			&n.Tick,
		)
		if err != nil {
			logs.Log.Error("Failed to scan news row at tick", zap.Int("tick", tick), zap.Error(err))
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}
