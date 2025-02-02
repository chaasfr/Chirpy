package main

import "net/http"

	const filepathRoot = "."
	const port = "8080"

func InitServer() {
	mux := http.NewServeMux()
	mux.Handle("/",http.FileServer(http.Dir(filepathRoot)))

	var httpServer http.Server
	httpServer.Handler = mux
	httpServer.Addr = ":" + port


	httpServer.ListenAndServe()
}