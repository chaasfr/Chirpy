package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chaasfr/chirpy/internal/auth"
)

type LoginInput struct{
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds *int64 `json:"expires_in_seconds,omitempty"`
}

type userWithTokenJson struct{
	userJson
	Token    string `json:"token"`
}

func (cfg *apiConfig)HandlerLogin(rw http.ResponseWriter, req *http.Request){
	var input LoginInput

	if err:=GetInputStruct(&input, rw, req); err != nil {
		return
	}

	userDb, err := cfg.dbQueries.GetUserPassword(req.Context(), input.Email)
	if  err != nil {
		if strings.Contains(err.Error(), "no row") {
			ReturnJsonError(rw, 401, "incorrect email or password")
			return
		} else {
			log.Printf("error retrieveing password %s", err)
			ReturnJsonGenericInternalError(rw)
			return
		}
	}

	err = auth.CheckPasswordHash(input.Password,userDb.HashedPassword)
	if err != nil {
		ReturnJsonError(rw, 401, "incorrect email or password")
		return
	}

	jwtDuration := auth.JwtDefaultDuration
	if input.ExpiresInSeconds != nil {
		jwtDuration = time.Duration(*input.ExpiresInSeconds) * time.Second
	}
	token, err := auth.MakeJWT(userDb.ID, cfg.jwtSecret, jwtDuration)

	if err != nil {
		log.Printf("error creating JWT token %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}
	userJson := userJsonFromDb(userDb)
	userWithTokenJson := userWithTokenJson{userJson: userJson, Token: token}
	ReturnWithJSON(rw, 200, userWithTokenJson)
}