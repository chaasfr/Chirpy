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
		ReturnGenericJsonError(rw, req)
		return
	}

	if len(input.Body) > 140 {
		ReturnTooLongChirpError(rw, req)
		return
	}

	ReturnValidChirp(rw, req, input.Body)
}

func ReturnValidChirp(rw http.ResponseWriter, req *http.Request, chirp string){
	badWords := []string {"kerfuffle", "sharbert", "fornax"}
	chirpSlice := strings.Split(chirp, " ")
	for i, word := range chirpSlice{
		if slices.Contains(badWords, strings.ToLower(word)) {
			chirpSlice[i] = "****"
		}
	}

	outputValid := ValidateChirpOutputValid{true, strings.Join(chirpSlice," ")}
	dat, err := json.Marshal(outputValid)
	if err != nil {
		log.Printf("error marshalling json %s",err)
		ReturnGenericJsonError(rw, req)
		return
	}
	rw.WriteHeader(200)
	rw.Write(dat)
}


func ReturnTooLongChirpError(rw http.ResponseWriter, req *http.Request){
	outputError := GenericJsonError{"Chirp is too long"}
	dat, err := json.Marshal(outputError)
	if err != nil {
		log.Printf("error marshalling json %s",err)
		rw.WriteHeader(500)
		return
	}
	rw.WriteHeader(400)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(dat)
}


func ReturnGenericJsonError(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(500)
	outputError := GenericJsonError{"Something went wrong"}
	dat, err := json.Marshal(outputError)
	if err != nil {
		log.Printf("error marshalling json %s",err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(dat)
}