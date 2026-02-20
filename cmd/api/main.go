package main

import (
	"context"
	"fmt"
	"go-auth/internal/app"
	"go-auth/internal/config"
	"go-auth/internal/server"
	"log"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Error occured loading config: %v", err)
	}


	// connect database from app

	a, err := app.New(context.Background())



	if err != nil {
		log.Fatalf("Error occured in database: %v", err)
	}
	
	

	defer func() {
		if err := a.Close(context.Background()); err != nil {
			log.Printf("Mongo Disconnected: %v", err)
		}
	}()



	//router

	router := server.NewRouter(a, cfg)
	

	add := fmt.Sprintf(":%s", cfg.PORT)

 if err := 	router.Run(add);err != nil {
	log.Fatalf("Server Error")
 }

}
