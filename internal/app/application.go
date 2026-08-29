package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Application struct {
	Config *Config
	Logger *zap.Logger

	DB    *pgxpool.Pool
	Redis *redis.Client

	Router *gin.Engine
}