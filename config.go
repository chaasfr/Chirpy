package main

import (
	"database/sql"
	"log"
	
	"os"
	"sync/atomic"

	"github.com/chaasfr/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	jwtSecret      string
}

func initConfig() *apiConfig {
	var cfg apiConfig
	cfg.dbQueries = initDbConnection()
	cfg.platform = os.Getenv("PLATFORM")
	cfg.jwtSecret = os.Getenv("JWTSECRET")
	return &cfg
}

func initDbConnection() *database.Queries {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Printf("error connecting to db %s", err)
		os.Exit(1)
	}

	return database.New(db)
}

