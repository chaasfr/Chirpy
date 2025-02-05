package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/chaasfr/chirpy/internal/auth"
)

type LoginInput struct{
	Password string `json:"password"`
	Email string `json:"email"`
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

	output := userJsonFromDb(userDb)
	ReturnWithJSON(rw, 200, output)
}