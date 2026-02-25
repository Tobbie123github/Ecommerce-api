package app

import (
	"context"
	"fmt"
	"go-auth/internal/config"
	"go-auth/internal/db"
	"time"

	//"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	Config config.Config

	MongoClient *mongo.Client
	DB 	*mongo.Database
	Redis *redis.Client
}



func New(ctx context.Context) (*App, error) {
	
	
	// Connect mongo and config

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("Cant load config: %v", err)
	}

	mongo, err := db.Connect(ctx, cfg)
	
	if err != nil {
		return nil, fmt.Errorf("Cant load db: %v", err)
	}
	redisCli, err := db.Redis(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("Cant load db: %v", err)
	}
	

	return &App{
		Config: cfg,
		MongoClient: mongo.Client,
		DB: mongo.DB,
		Redis: redisCli,
	}, nil


}

func (a *App) Close(ctx context.Context) error {

	if a.MongoClient == nil {
		return nil
	}

	closeCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	return a.MongoClient.Disconnect(closeCtx)


}