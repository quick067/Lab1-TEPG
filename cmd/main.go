package main

import (
	"log"
	"training-system/internal/config"
	"training-system/internal/server"
	"training-system/internal/storage"
)

func main(){
	cfg := config.Load()

	db, err := storage.NewDBConnection(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}

	server := server.NewServer(db, *cfg)
	if err := server.RunServer(); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}

//docker run --name training-container --env-file .env training-app
//docker run --name training-container --env-file .env -p 8080:8080 training-app