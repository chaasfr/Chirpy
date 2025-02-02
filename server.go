package main

import (
	"net/http"
)

	const filepathRoot = "."
	const readynessEndpoint = "GET /api/healthz"
	const fileServerEndpoint = "/app/"
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

	var httpServer http.Server
	httpServer.Handler = mux
	httpServer.Addr = ":" + port


	httpServer.ListenAndServe()
}