package main

import (
	"fmt"
	"net/http"
)

	const filepathRoot = "."
	const readynessEndpoint = "GET /healthz"
	const fileServerEndpoint = "/app/"
	const metricsEndpoint = "GET /metrics"
	const resetMetricsEndpoint ="POST /reset"

	const port = "8080"

func InitServer() {
	mux := http.NewServeMux()
	cfg := initConfig()
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle(fileServerEndpoint, cfg.middlewareMetricsInc(fileServer))
	mux.HandleFunc(readynessEndpoint, HandleReadinessReq)
	mux.HandleFunc(metricsEndpoint, cfg.HanddleMetricsReq)
	mux.HandleFunc(resetMetricsEndpoint, cfg.HanddleResetReq)

	var httpServer http.Server
	httpServer.Handler = mux
	httpServer.Addr = ":" + port


	httpServer.ListenAndServe()
}

func HandleReadinessReq(rw http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte("OK"))
}

func (cfg *apiConfig) HanddleMetricsReq(rw http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	text := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	rw.Write([]byte(text))
}

func (cfg *apiConfig) HanddleResetReq(rw http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Swap(0)
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte("Reset of metrics executed"))
}