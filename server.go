package main

import (
	"net/http"
)

const filepathRoot = "."

const readynessEndpoint = "GET /api/healthz"
const fileServerEndpoint = "/app/"
const CreateUserEndpoint = "POST /api/users"
const CreateChirpEndpoint = "POST /api/chirps"
const GetChirpsEndpoint = "GET /api/chirps"
const GetChirpByIDEndpoint = "GET /api/chirps/{chirpID}"
const LoginEndpoint = "POST /api/login"
const RefreshEndpoint = "POST /api/refresh"
const RevokeEndpoint = "POST /api/revoke"

const metricsEndpoint = "GET /admin/metrics"
const resetMetricsEndpoint = "POST /admin/reset"

const port = "8080"

func InitServer() {
	mux := http.NewServeMux()
	cfg := initConfig()
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle(fileServerEndpoint, cfg.mdwMetricsInc(fileServer))
	mux.HandleFunc(readynessEndpoint, HandlerReadiness)
	mux.HandleFunc(metricsEndpoint, cfg.HandlerMetrics)
	mux.HandleFunc(resetMetricsEndpoint, cfg.mdwDevPlatformOnly(cfg.HandlerReset))
	mux.HandleFunc(CreateChirpEndpoint, cfg.mdwValidateJWT(cfg.HandlerCreateChirp))
	mux.HandleFunc(CreateUserEndpoint, cfg.HandlerCreateUser)
	mux.HandleFunc(GetChirpsEndpoint, cfg.HandlerGetChirps)
	mux.HandleFunc(GetChirpByIDEndpoint, cfg.HandlerGetChirpById)
	mux.HandleFunc(LoginEndpoint, cfg.HandlerLogin)
	mux.HandleFunc(RefreshEndpoint, cfg.HandlerRefresh)
	mux.HandleFunc(RevokeEndpoint, cfg.HandlerRevoke)

	var httpServer http.Server
	httpServer.Handler = mux
	httpServer.Addr = ":" + port

	httpServer.ListenAndServe()
}
