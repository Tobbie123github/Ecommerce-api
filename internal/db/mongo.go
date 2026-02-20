package db

import (
	"context"
	"fmt"
	"go-auth/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)


type Mongo struct {
	Client *mongo.Client
	DB *mongo.Database
}


func Connect(ctx context.Context, cfg config.Config) (*Mongo, error){

	chidCtx, cancel :=	context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	clientOpt := options.Client().ApplyURI(cfg.MONGO_URI)

	client, err := mongo.Connect(clientOpt)

	if err != nil {
		return nil, fmt.Errorf("Error connecting to mongo client: %v", err)

	}

	if err := client.Ping(chidCtx, nil); err != nil {
		return nil, fmt.Errorf("Mongo Ping Failed: %v", err)

	}


	// connect to mongodb

	database := client.Database(cfg.MONGO_DB)

	return &Mongo{
		Client: client,
		DB: database,
	}, nil

}