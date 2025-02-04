package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) HandlerReset(rw http.ResponseWriter, req *http.Request) {
	if err :=cfg.dbQueries.DeleteUsers(req.Context()); err != nil {
		log.Printf("error deleting users %s", err)
		ReturnGenericJsonError(rw)
		return
	}
	cfg.fileserverHits.Swap(0)
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte("Reset of metrics executed"))
}