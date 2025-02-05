package main

import (
	"log"
	"net/http"
)

type GetChirpsOutput []ChirpJson

func (cfg *apiConfig)HandlerGetChirps(rw http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.dbQueries.GetAllChirp(req.Context())

	if err != nil {
		log.Printf("error retrieving chirps fro, db %s", err)
		ReturnGenericJsonError(rw)
		return
	}

	output := GetChirpsOutput{}

	for _, chirpDb := range dbChirps{
		chirpJson := ChirpJsonFromDb(&chirpDb)
		output = append(output, *chirpJson)
	}

	ReturnWithJSON(rw, 200, output)
}