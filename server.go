package main

import (
	"net/http"
)

	const filepathRoot = "."

	const readynessEndpoint = "GET /api/healthz"
	const fileServerEndpoint = "/app/"
	const validateChirpEndpoint = "POST /api/validate_chirp"
	const CreateUserEndpoint = "POST /api/users"

	const metricsEndpoint = "GET /admin/metrics"
	const resetMetricsEndpoint ="POST /admin/reset"

	const port = "8080"

func InitServer() {
	mux := http.NewServeMux()
	cfg := initConfig()
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle(fileServerEndpoint, cfg.middlewareMetricsInc(fileServer))
	mux.HandleFunc(readynessEndpoint, HandlerReadiness)
	mux.HandleFunc(metricsEndpoint, cfg.HandlerMetrics)
	mux.HandleFunc(resetMetricsEndpoint, cfg.HandlerReset)
	mux.HandleFunc(validateChirpEndpoint, HandlerValidateChirpReq)
	mux.HandleFunc(CreateUserEndpoint, cfg.HandlerCreateUser)

	var httpServer http.Server
	httpServer.Handler = mux
	httpServer.Addr = ":" + port


	httpServer.ListenAndServe()
}