package main

import (
	"net/http"
	"time"

	"github.com/chaasfr/chirpy/internal/auth"
)

type RefreshJson struct{
	Token string `json:"token"`
}

func (cfg *apiConfig) HandlerRefresh(rw http.ResponseWriter, req *http.Request) {
	rtoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		ReturnJsonError(rw, 401, "no refresh token in authorization header")
		return
	}

	rtokenDb, err := cfg.dbQueries.GetRefreshToken(req.Context(), rtoken)
	if err != nil {
		ReturnJsonError(rw, 401, "unknown refresh token")
		return
	}

	if time.Now().After(rtokenDb.ExpiresAt) {
		ReturnJsonError(rw, 401, "refresh token expired")
		return
	}

	if rtokenDb.RevokedAt.Valid {
		ReturnJsonError(rw, 401, "refresh token has been revoked")
		return
	}

	jwt, err := auth.MakeJWT(rtokenDb.UserID, cfg.jwtSecret, auth.JwtDefaultDuration)
	if err != nil {
		ReturnJsonGenericInternalError(rw)
		return
	}

	refreshJson := RefreshJson{Token:jwt}
	ReturnWithJSON(rw, 200, refreshJson)
}
