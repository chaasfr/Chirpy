package main

import (
	"net/http"
	"sync/atomic"
	"database/sql"
	"log"
	"os"

	"github.com/chaasfr/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform string
}

func initConfig() *apiConfig{
	var cfg apiConfig
	cfg.dbQueries = initDbConnection()
	cfg.platform = os.Getenv("PLATFORM")
	return &cfg
}

func initDbConnection() *database.Queries{
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Printf("error connecting to db %s", err)
		os.Exit(1)
	}

	return database.New(db)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc( func(rw http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(rw, req)
	})
}

func (cfg *apiConfig) middlewareDevPlatformOnly(next func (http.ResponseWriter, *http.Request)) func (http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if cfg.platform != "dev" {
			log.Printf("This endpoint is only for dev %s", req.RequestURI)
			ReturnJsonError(rw, 403, "forbidden outside of dev platform")
			return
		}
		next(rw, req)
	}
}