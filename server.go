package main

import (
	"net/http"
)

const FilepathRoot = "."

const ReadynessEndpoint = "GET /api/healthz"
const FileServerEndpoint = "/app/"

const CreateUserEndpoint = "POST /api/users"
const UpdateUserEndpoint = "PUT /api/users"

const CreateChirpEndpoint = "POST /api/chirps"
const GetChirpsEndpoint = "GET /api/chirps"
const GetChirpByIDEndpoint = "GET /api/chirps/{" + chirpIdPathValue + "}"
const DeleteChirpByIdEndpoint = "DELETE /api/chirps/{" + chirpIdPathValue + "}"

const LoginEndpoint = "POST /api/login"
const RefreshEndpoint = "POST /api/refresh"
const RevokeEndpoint = "POST /api/revoke"

const MetricsEndpoint = "GET /admin/metrics"
const ResetMetricsEndpoint = "POST /admin/reset"

const PolkaWebhookEndpoint = "POST /api/polka/webhooks"

const port = "8080"

func InitServer() {
	mux := http.NewServeMux()
	cfg := initConfig()
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(FilepathRoot)))
	mux.HandleFunc(ReadynessEndpoint, HandlerReadiness)
	mux.Handle(FileServerEndpoint, cfg.mdwMetricsInc(fileServer))
	
	mux.HandleFunc(CreateUserEndpoint, cfg.HandlerCreateUser)
	mux.HandleFunc(UpdateUserEndpoint, cfg.mdwValidateJWT(cfg.HandlerUpdateUser))

	
	mux.HandleFunc(CreateChirpEndpoint, cfg.mdwValidateJWT(cfg.HandlerCreateChirp))
	
	mux.HandleFunc(GetChirpsEndpoint, cfg.HandlerGetChirps)
	mux.HandleFunc(GetChirpByIDEndpoint, cfg.HandlerGetChirpById)
	mux.HandleFunc(DeleteChirpByIdEndpoint, cfg.mdwValidateJWT(cfg.HandlerDeleteChirpById))

	mux.HandleFunc(LoginEndpoint, cfg.HandlerLogin)
	mux.HandleFunc(RefreshEndpoint, cfg.HandlerRefresh)
	mux.HandleFunc(RevokeEndpoint, cfg.HandlerRevoke)

	mux.HandleFunc(MetricsEndpoint, cfg.HandlerMetrics)
	mux.HandleFunc(ResetMetricsEndpoint, cfg.mdwDevPlatformOnly(cfg.HandlerReset))

	mux.HandleFunc(PolkaWebhookEndpoint, cfg.HandlerPolkaWebhook)

	var httpServer http.Server
	httpServer.Handler = mux
	httpServer.Addr = ":" + port

	httpServer.ListenAndServe()
}
