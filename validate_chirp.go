package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

type ValidateChirpInput struct {
	Body string `json:"body"`
}

type GenericJsonError struct {
	Error string `json:"error"`
}

type ValidateChirpOutputValid struct {
	Valid bool `json:"valid"`
	Cleaned_body string `json:"cleaned_body"`
}

func HandlerValidateChirpReq(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/json")

	input := ValidateChirpInput{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&input)
	if err != nil {
		log.Printf("Error decoding Chirp to valid %s", err)	
		ReturnGenericJsonError(rw)
		return
	}

	if len(input.Body) > 140 {
		ReturnJsonError(rw, 400, "Chirp is too long")
		return
	}

	ReturnValidChirp(rw, input.Body)
}

func ReturnValidChirp(rw http.ResponseWriter, chirp string){
	badWords := []string {"kerfuffle", "sharbert", "fornax"}
	chirpSlice := strings.Split(chirp, " ")
	for i, word := range chirpSlice{
		if slices.Contains(badWords, strings.ToLower(word)) {
			chirpSlice[i] = "****"
		}
	}
	outputValid := ValidateChirpOutputValid{true, strings.Join(chirpSlice," ")}
	ReturnWithJSON(rw, 200, outputValid)
}