package main

import (
	"context"
	"log"
	"net/http"

	"github.com/chaasfr/chirpy/internal/auth"
)

const UseridFromJwtKey = "user_id_from_jwt"

func (cfg *apiConfig) mdwMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(rw, req)
	})
}

func (cfg *apiConfig) mdwDevPlatformOnly(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if cfg.platform != "dev" {
			log.Printf("This endpoint is only for dev %s", req.RequestURI)
			ReturnJsonError(rw, 403, "forbidden outside of dev platform")
			return
		}
		next(rw, req)
	}
}

func (cfg *apiConfig) mdwValidateJWT(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		token, err := auth.GetBearerToken(req.Header)
		if err != nil {
			ReturnJsonError(rw, 401, "no bearer token in authorization header")
			return
		}
		uuid, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			ReturnJsonError(rw, 401, "invalid or expired token")
			return
		}
		
		ctxWithId := context.WithValue(req.Context(), UseridFromJwtKey, uuid)
		newReq := req.WithContext(ctxWithId)

		next(rw, newReq)
	}
}