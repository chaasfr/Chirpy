package main

import (
	"net/http"
)

const filepathRoot = "."

const readynessEndpoint = "GET /api/healthz"
const fileServerEndpoint = "/app/"

const CreateUserEndpoint = "POST /api/users"
const UpdateUserEndpoint = "PUT /api/users"

const CreateChirpEndpoint = "POST /api/chirps"
const GetChirpsEndpoint = "GET /api/chirps"
const GetChirpByIDEndpoint = "GET /api/chirps/{" + chirpIdPathValue + "}"
const DeleteChirpByIdEndpoint = "DELETE /api/chirps/{" + chirpIdPathValue + "}"

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
	mux.HandleFunc(readynessEndpoint, HandlerReadiness)
	mux.Handle(fileServerEndpoint, cfg.mdwMetricsInc(fileServer))
	
	mux.HandleFunc(CreateUserEndpoint, cfg.HandlerCreateUser)
	mux.HandleFunc(UpdateUserEndpoint, cfg.mdwValidateJWT(cfg.HandlerUpdateUser))

	
	mux.HandleFunc(CreateChirpEndpoint, cfg.mdwValidateJWT(cfg.HandlerCreateChirp))
	
	mux.HandleFunc(GetChirpsEndpoint, cfg.HandlerGetChirps)
	mux.HandleFunc(GetChirpByIDEndpoint, cfg.HandlerGetChirpById)
	mux.HandleFunc(DeleteChirpByIdEndpoint, cfg.mdwValidateJWT(cfg.HandlerDeleteChirpById))

	mux.HandleFunc(LoginEndpoint, cfg.HandlerLogin)
	mux.HandleFunc(RefreshEndpoint, cfg.HandlerRefresh)
	mux.HandleFunc(RevokeEndpoint, cfg.HandlerRevoke)

	mux.HandleFunc(metricsEndpoint, cfg.HandlerMetrics)
	mux.HandleFunc(resetMetricsEndpoint, cfg.mdwDevPlatformOnly(cfg.HandlerReset))

	var httpServer http.Server
	httpServer.Handler = mux
	httpServer.Addr = ":" + port

	httpServer.ListenAndServe()
}
