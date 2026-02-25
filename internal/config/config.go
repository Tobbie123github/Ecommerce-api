package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MONGO_URI string
	MONGO_DB  string
	PORT      string
	JWTSecret string
	CLOUDINARY_URL string
	REDIS_URL string
	STRIPE_SECRET_KEY string
	STRIPE_WEBHOOK_SECRET string
}

func Load() (Config, error) {

	err := godotenv.Load()

	if err != nil {
		return Config{}, err
	}

	mongoURI, err := extractText("MONGO_URI")

	if err != nil {
		return Config{}, fmt.Errorf("MongoURI not found") 
	}

	mongoDB, err := extractText("MONGO_DB")

	if err != nil {
		return Config{}, fmt.Errorf("MONGO_DB not found") 
	}

	port, err := extractText("PORT")

	if err != nil {
		return Config{}, fmt.Errorf("PORT not found") 
	}

	jwtSecret, err := extractText("JWTSecret")

	if err != nil {
		return Config{}, fmt.Errorf("jwtSecret not found") 
	}

	cloudinary, err := extractText("CLOUDINARY_URL")

	if err != nil {
		return Config{}, fmt.Errorf("cloudinary not found") 
	}

	redisURL, err := extractText("REDIS_URL")

	if err != nil {
		return Config{}, fmt.Errorf("redis url not found") 
	}

	stripe, err := extractText("STRIPE_SECRET_KEY")

	if err != nil {
		return Config{}, fmt.Errorf("stripe secrete key not found") 
	}

	webhook, err := extractText("STRIPE_WEBHOOK_SECRET")

	if err != nil {
		return Config{}, fmt.Errorf("stripe webhook secrete not found") 
	}
	
	return Config{
		MONGO_URI: mongoURI,
		MONGO_DB: mongoDB,
		PORT: port,
		JWTSecret: jwtSecret,
		CLOUDINARY_URL: cloudinary,
		REDIS_URL: redisURL,
		STRIPE_SECRET_KEY: stripe,
		STRIPE_WEBHOOK_SECRET: webhook,
	}, nil

}

func extractText(key string) (string, error) {

	value := os.Getenv(key)

	if value == "" {
		return "", fmt.Errorf("Env cannot have an empty value")
	}

	return value, nil


}