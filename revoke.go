package main

import (
	"net/http"

	"github.com/chaasfr/chirpy/internal/auth"
)

func (cfg *apiConfig) HandlerRevoke(rw http.ResponseWriter, req *http.Request) {
	rtoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		ReturnJsonError(rw, 401, "no refresh token in authorization header")
		return
	}

	err = cfg.DbQueries.RevokeRefreshToken(req.Context(), rtoken)
	if err != nil {
		ReturnJsonError(rw, 401, "unknown refresh token")
		return
	}

	rw.WriteHeader(204)
}
