package main

import "net/http"

func (cfg *apiConfig) HandlerReset(rw http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Swap(0)
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte("Reset of metrics executed"))
}