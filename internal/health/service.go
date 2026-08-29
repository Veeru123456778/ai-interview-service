package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	db      *pgxpool.Pool
	redis   *redis.Client
	appName string
	version string
}

func NewService(
	db *pgxpool.Pool,
	redis *redis.Client,
	appName string,
	version string,
) *Service {
	return &Service{
		db:      db,
		redis:   redis,
		appName: appName,
		version: version,
	}
}

func (s *Service) Live() map[string]any {
	return map[string]any{
		"status": "ALIVE",
	}
}

func (s *Service) Health() map[string]any {
	return map[string]any{
		"status":  "UP",
		"service": s.appName,
		"version": s.version,
	}
}

func (s *Service) Ready(ctx context.Context) map[string]any {
	postgres := "UP"
	redis := "UP"
	status := "READY"

	if err := s.db.Ping(ctx); err != nil {
		status = "NOT_READY"
		postgres = "DOWN"
	}

	if err := s.redis.Ping(ctx).Err(); err != nil {
		status = "NOT_READY"
		redis = "DOWN"
	}

	return map[string]any{
		"status": status,
		"dependencies": map[string]string{
			"postgres": postgres,
			"redis":    redis,
		},
	}
}