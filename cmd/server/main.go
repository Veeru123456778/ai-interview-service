package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Veeru123456778/ai-interview-service/internal/app"
	"github.com/Veeru123456778/ai-interview-service/internal/storage"
	"github.com/Veeru123456778/ai-interview-service/internal/shared/logger"
)

func main() {

	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(cfg.App.Environment)
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewPostgres(context.Background(), cfg.Database.URL)
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := storage.NewRedis(context.Background(), cfg.Redis.URL)
	if err != nil {
		log.Fatal(err)
	}

	application := &app.Application{
		Config: cfg,
		Logger: logg,
		DB:     db,
		Redis:  redisClient,
	}

	router := app.NewRouter(application)
	application.Router = router

	address := fmt.Sprintf(":%d", cfg.Server.Port)

	logg.Info("starting server", zap.String("address", address))

	if err := router.Run(address); err != nil {
		log.Fatal(err)
	}
}