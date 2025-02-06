package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/chaasfr/chirpy/internal/auth"
	"github.com/chaasfr/chirpy/internal/database"
)

type LoginInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userWithTokenJson struct {
	userJson
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) HandlerLogin(rw http.ResponseWriter, req *http.Request) {
	var input LoginInput

	if err := GetInputStructFromJson(&input, rw, req); err != nil {
		return
	}

	userDb, err := cfg.DbQueries.GetUserPassword(req.Context(), input.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no row") {
			ReturnJsonError(rw, 401, "incorrect email or password")
			return
		} else {
			log.Printf("error retrieveing password %s", err)
			ReturnJsonGenericInternalError(rw)
			return
		}
	}

	err = auth.CheckPasswordHash(input.Password, userDb.HashedPassword)
	if err != nil {
		ReturnJsonError(rw, 401, "incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(userDb.ID, cfg.JWTSecret, auth.JwtDefaultDuration)
	if err != nil {
		log.Printf("error creating JWT token %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("error creating refresh token %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	createRefreshTokenQp := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: userDb.ID,
	}
	_, err = cfg.DbQueries.CreateRefreshToken(req.Context(), createRefreshTokenQp)
	if err != nil {
		log.Printf("error storing refresh token %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	userJson := userJsonFromDb(userDb)
	userWithTokenJson := userWithTokenJson{
		userJson:     userJson,
		Token:        token,
		RefreshToken: refreshToken,
	}
	ReturnWithJSON(rw, 200, userWithTokenJson)
}
