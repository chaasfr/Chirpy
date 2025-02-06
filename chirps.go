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

const chirpIdPathValue = "chirpID"

type CreateChirpInput struct {
	Body string `json:"body"`
}

type ChirpJson struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type GetChirpsOutput []ChirpJson

func ChirpJsonFromDb(chirp *database.Chirp) *ChirpJson {
	return &ChirpJson{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
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

func (cfg *apiConfig) HandlerCreateChirp(rw http.ResponseWriter, req *http.Request) {

	input := CreateChirpInput{}
	if err := GetInputStructFromJson(&input, rw, req); err != nil {
		return
	}

	if len(input.Body) > 140 {
		ReturnJsonError(rw, 400, "Chirp is too long")
		return
	}

	chirpBody := validateChirp(input.Body)
	userId, ok := req.Context().Value(UseridFromJwtKey).(uuid.UUID)
	if !ok {
		log.Printf("Error retrieveing userId from ctx\n")
		ReturnJsonGenericInternalError(rw)
		return
	}

	qp := database.CreateChirpParams{Body: chirpBody, UserID: userId}
	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), qp)
	if err != nil {
		log.Printf("error saving chirp to db %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	output := ChirpJsonFromDb(&chirp)

	ReturnWithJSON(rw, 201, output)
}

func (cfg *apiConfig) HandlerGetChirps(rw http.ResponseWriter, req *http.Request) {
	chirpsDb, err := cfg.dbQueries.GetAllChirp(req.Context())

	if err != nil {
		log.Printf("error retrieving chirps fro, db %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	output := GetChirpsOutput{}

	for _, chirpDb := range chirpsDb {
		chirpJson := ChirpJsonFromDb(&chirpDb)
		output = append(output, *chirpJson)
	}

	ReturnWithJSON(rw, 200, output)
}

func (cfg *apiConfig) HandlerGetChirpById(rw http.ResponseWriter, req *http.Request) {
	chirpIdInput := req.PathValue(chirpIdPathValue)

	chirpId, err := uuid.Parse(chirpIdInput)
	if err != nil {
		log.Printf("error parsing chirpId %s", err)
		ReturnJsonError(rw, 404, "chirp not found - id is invalid")
		return
	}
	chirpDb, err := cfg.dbQueries.GetChirp(req.Context(), chirpId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			ReturnJsonError(rw, 404, "chirp not found")
			return
		} else {
			log.Printf("error retrieving chirp %s", err)
			ReturnJsonGenericInternalError(rw)
			return
		}

	}

	chirpJson := ChirpJsonFromDb(&chirpDb)

	ReturnWithJSON(rw, 200, chirpJson)
}

func (cfg *apiConfig) HandlerDeleteChirpById(rw http.ResponseWriter, req *http.Request) {
	chirpIdText := req.PathValue(chirpIdPathValue)

	chirpId, err := uuid.Parse(chirpIdText)
	if err != nil {
		ReturnJsonError(rw, 404, "chirpId is not a valid uuid")
		return
	}

	chirpDb, err := cfg.dbQueries.GetChirp(req.Context(), chirpId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			ReturnJsonError(rw, 404, "chirp not found")
			return
		} else {
			log.Printf("error retrieving chirp %s", err)
			ReturnJsonGenericInternalError(rw)
			return
		}
	}

	userId, ok := req.Context().Value(UseridFromJwtKey).(uuid.UUID)
	if !ok {
		log.Printf("Error retrieveing userId from ctx\n")
		ReturnJsonGenericInternalError(rw)
		return
	}

	if chirpDb.UserID != userId {
		ReturnJsonError(rw, 403, "unauthorized - user not associated to chirp")
		return
	}

	err = cfg.dbQueries.DeleteChirp(req.Context(), chirpId)
	if err != nil {
		log.Printf("error deleting from DB chirp %v", chirpId)
		ReturnJsonGenericInternalError(rw)
		return
	}

	rw.WriteHeader(204)
}
