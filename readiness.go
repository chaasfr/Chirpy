package main

import "net/http"

func HandlerReadiness(rw http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write([]byte("OK"))
}