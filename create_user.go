package main

import (
	"log"
	"net/http"

	"github.com/chaasfr/chirpy/internal/auth"
	"github.com/chaasfr/chirpy/internal/database"
)

type CreateUserInput struct {
	Password string `json:"password"`
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

	hashed_password, err := auth.HashPassword(input.Password)
	if err != nil {
		log.Printf("Error hashing password %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	qp := database.CreateUserParams{Email: input.Email, HashedPassword: hashed_password}

	user, err := cfg.dbQueries.CreateUser(req.Context(), qp)
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
