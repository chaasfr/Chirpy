package main

import (
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/chaasfr/chirpy/internal/database"
	"github.com/google/uuid"
)

type CreateChirpInput struct {
	Body   string `json:"body"`
	UserId string `json:"user_id"`
}

type ChirpJson struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func ChirpJsonFromDb(chirp *database.Chirp) *ChirpJson {
	return &ChirpJson{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
}

func (cfg *apiConfig) HandlerCreateChirp(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/json")

	input := CreateChirpInput{}
	if err:=GetInputStruct(&input, rw, req); err != nil {
		return
	}

	if len(input.Body) > 140 {
		ReturnJsonError(rw, 400, "Chirp is too long")
		return
	}

	chirpBody := validateChirp(input.Body)
	uuidId, err := uuid.Parse(input.UserId)
	if err != nil {
		log.Printf("error converting uuid %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	qp := database.CreateChirpParams{Body: chirpBody, UserID: uuidId}
	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), qp)
	if err != nil {
		log.Printf("error saving chirp to db %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	output := ChirpJsonFromDb(&chirp)

	ReturnWithJSON(rw, 201, output)
}

func validateChirp(chirp string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpSlice := strings.Split(chirp, " ")
	for i, word := range chirpSlice {
		if slices.Contains(badWords, strings.ToLower(word)) {
			chirpSlice[i] = "****"
		}
	}
	return strings.Join(chirpSlice, " ")
}
