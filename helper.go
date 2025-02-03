package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func ReturnGenericJsonError(rw http.ResponseWriter) {
	ReturnJsonError(rw, 500, "something went wrong")
}

func ReturnJsonError(rw http.ResponseWriter, code int, msg string){
	outputError := GenericJsonError{msg}
	dat, err := json.Marshal(outputError)
	if err != nil {
		log.Printf("error marshalling json %s",err)
		rw.WriteHeader(500)
		return
	}
	rw.WriteHeader(code)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(dat)
}

func ReturnWithJSON(rw http.ResponseWriter, code int, payload interface{}){
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling json %s",err)
		ReturnGenericJsonError(rw)
		return
	}
	rw.WriteHeader(code)
	rw.Write(dat)
}