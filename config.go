package main

import (
	"database/sql"
	"log"

	"os"
	"sync/atomic"

	"github.com/chaasfr/chirpy/internal/database"
)

type apiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
	JWTSecret      string
	PolkaKey       string
}

func initConfig() *apiConfig {
	var cfg apiConfig
	cfg.DbQueries = initDbConnection()
	cfg.Platform = os.Getenv("PLATFORM")
	cfg.JWTSecret = os.Getenv("JWTSECRET")
	cfg.PolkaKey = os .Getenv("POLKA_KEY")
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
