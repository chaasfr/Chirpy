package main

import "net/http"


func InitServer() {
	serverMux := http.NewServeMux()

	var httpServer http.Server

	httpServer.Handler = serverMux
	httpServer.Addr = ":8080"

	httpServer.ListenAndServe()
}