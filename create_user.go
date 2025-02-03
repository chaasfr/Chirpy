package main

import (
	"encoding/json"
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


func (cfg *apiConfig)HandlerCreateUser(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/json")

	input := CreateUserInput{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&input)
	if err != nil {
		log.Printf("Error decoding Create User %s", err)	
		ReturnGenericJsonError(rw)
		return
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(),input.Email)
	if err != nil {
		log.Printf("Error creating user %s", err)
		ReturnGenericJsonError(rw)
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