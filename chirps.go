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

const ChirpIdPathValue = "chirpID"
const AuthorIdQueryParamKey = "author_id"
const SortQueryParamKey = "sort"
const SortQueryDescKeyword = "desc"

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
	chirp, err := cfg.DbQueries.CreateChirp(req.Context(), qp)
	if err != nil {
		log.Printf("error saving chirp to db %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	output := ChirpJsonFromDb(&chirp)

	ReturnWithJSON(rw, 201, output)
}

func sortChips(a, b database.Chirp) bool {
	return a.CreatedAt.After(b.CreatedAt)
}

func (cfg *apiConfig) HandlerGetChirps(rw http.ResponseWriter, req *http.Request) {
	var chirpsDb []database.Chirp
	var err error

	if req.URL.Query().Has(AuthorIdQueryParamKey) {
		authorIdStr := req.URL.Query().Get(AuthorIdQueryParamKey)
		auhorId, errParse := uuid.Parse(authorIdStr)
		if errParse != nil {
			ReturnJsonError(rw, 500, "invalid " + AuthorIdQueryParamKey)
			return
		}
		chirpsDb, err = cfg.DbQueries.GetAllChirpFromUser(req.Context(), auhorId)
	} else {
		chirpsDb, err = cfg.DbQueries.GetAllChirp(req.Context())
	}

	if err != nil {
		log.Printf("error retrieving chirps from db %s", err)
		ReturnJsonGenericInternalError(rw)
		return
	}

	//sort by asc by default
	if req.URL.Query().Get(SortQueryParamKey) == SortQueryDescKeyword {
		slices.Reverse(chirpsDb)
	}
	
	output := GetChirpsOutput{}

	for _, chirpDb := range chirpsDb {
		chirpJson := ChirpJsonFromDb(&chirpDb)
		output = append(output, *chirpJson)
	}

	ReturnWithJSON(rw, 200, output)
}

func (cfg *apiConfig) HandlerGetChirpById(rw http.ResponseWriter, req *http.Request) {
	chirpIdInput := req.PathValue(ChirpIdPathValue)

	chirpId, err := uuid.Parse(chirpIdInput)
	if err != nil {
		log.Printf("error parsing chirpId %s", err)
		ReturnJsonError(rw, 404, "chirp not found - id is invalid")
		return
	}
	chirpDb, err := cfg.DbQueries.GetChirp(req.Context(), chirpId)
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
	chirpIdText := req.PathValue(ChirpIdPathValue)

	chirpId, err := uuid.Parse(chirpIdText)
	if err != nil {
		ReturnJsonError(rw, 404, "chirpId is not a valid uuid")
		return
	}

	chirpDb, err := cfg.DbQueries.GetChirp(req.Context(), chirpId)
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

	err = cfg.DbQueries.DeleteChirp(req.Context(), chirpId)
	if err != nil {
		log.Printf("error deleting from DB chirp %v", chirpId)
		ReturnJsonGenericInternalError(rw)
		return
	}

	rw.WriteHeader(204)
}
