package main

import (
	"log"
	"net/http"
)

type CreateUserInput struct {
	Email string `json:"email"`
}

type CreateUserOutput struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}

func (cfg *apiConfig) HandlerCreateUser(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/json")

	input := CreateUserInput{}
	GetInputStruct(&input, rw, req)

	user, err := cfg.dbQueries.CreateUser(req.Context(), input.Email)
	if err != nil {
		log.Printf("Error creating user %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	output := CreateUserOutput{
		user.ID.String(),
		user.CreatedAt.String(),
		user.UpdatedAt.String(),
		user.Email,
	}
	ReturnWithJSON(rw, 201, output)
}
