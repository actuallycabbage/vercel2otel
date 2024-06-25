package main

import (
	"net/http"
	handler "vercel2otel/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Handler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic("could not start server")
	}
}
