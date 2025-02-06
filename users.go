package main

import (
	"log"
	"net/http"

	"github.com/chaasfr/chirpy/internal/auth"
	"github.com/chaasfr/chirpy/internal/database"
	"github.com/google/uuid"
)

type UserInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userJson struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func userJsonFromDb(user database.User) userJson {
	return userJson{
		user.ID.String(),
		user.CreatedAt.String(),
		user.UpdatedAt.String(),
		user.Email,
		user.IsChirpyRed,
	}
}

func (cfg *apiConfig) HandlerCreateUser(rw http.ResponseWriter, req *http.Request) {

	input := UserInput{}
	if err := GetInputStructFromJson(&input, rw, req); err != nil {
		return
	}

	hashed_password, err := auth.HashPassword(input.Password)
	if err != nil {
		log.Printf("Error hashing password %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	qp := database.CreateUserParams{Email: input.Email, HashedPassword: hashed_password}

	user, err := cfg.DbQueries.CreateUser(req.Context(), qp)
	if err != nil {
		log.Printf("Error creating user %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	output := userJsonFromDb(user)
	ReturnWithJSON(rw, 201, output)
}

func (cfg *apiConfig) HandlerUpdateUser(rw http.ResponseWriter, req *http.Request) {
	input := UserInput{}
	if err := GetInputStructFromJson(&input, rw, req); err != nil {
		return
	}

	hashed_password, err := auth.HashPassword(input.Password)
	if err != nil {
		log.Printf("Error hashing password %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	userId, ok := req.Context().Value(UseridFromJwtKey).(uuid.UUID)
	if !ok {
		log.Printf("Error retrieveing userId from ctx\n")
		ReturnJsonGenericInternalError(rw)
		return
	}

	qp := database.UpdateUserParams{
		Email:          input.Email,
		HashedPassword: hashed_password,
		ID:             userId,
	}
	user, err := cfg.DbQueries.UpdateUser(req.Context(), qp)
	if err != nil {
		log.Printf("error updating user in db %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	userJson := userJsonFromDb(user)
	ReturnWithJSON(rw, 200, userJson)
}
