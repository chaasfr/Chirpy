package main

import (
	"log"
	"net/http"
	"time"

	"github.com/chaasfr/chirpy/internal/auth"
)

type RefreshJson struct{
	Token string `json:"token"`
}

func (cfg *apiConfig) HandlerRefresh(rw http.ResponseWriter, req *http.Request) {
	log.Println("refreshing...")
	rtoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		ReturnJsonError(rw, 401, "no refresh token in authorization header")
		return
	}
	log.Printf("got resfresh token %s\n", rtoken)

	rtokenDb, err := cfg.dbQueries.GetRefreshToken(req.Context(), rtoken)
	if err != nil {
		ReturnJsonError(rw, 401, "unknown refresh token")
		return
	}

	log.Printf("got resfresh token %s in DB\n", rtoken)

	if time.Now().After(rtokenDb.ExpiresAt) {
		ReturnJsonError(rw, 401, "refresh token expired")
		return
	}

	log.Printf("token %s is not expired \n", rtoken)

	jwt, err := auth.MakeJWT(rtokenDb.UserID, cfg.jwtSecret, auth.JwtDefaultDuration)
	if err != nil {
		ReturnJsonGenericInternalError(rw)
		return
	}

	log.Printf("created new jwt\n")

	refreshJson := RefreshJson{Token:jwt}
	ReturnWithJSON(rw, 200, refreshJson)
}
