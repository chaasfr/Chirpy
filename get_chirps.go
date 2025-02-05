package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type GetChirpsOutput []ChirpJson

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
	chirpIdInput := req.PathValue("chirpID")

	chirpId, err := uuid.Parse(chirpIdInput)
	if err != nil {
		log.Printf("error parsing chirpId %s", err)
		ReturnJsonError(rw, 404, "chirp not found - id is invalid")
		return
	}
	chirpDb, err := cfg.dbQueries.GetChirp(req.Context(), chirpId)
	if err != nil {
		if strings.Contains(err.Error(),"no rows") {
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
