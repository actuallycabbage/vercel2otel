package main

import (
	"context"
	"net/http"
)

func main() {
	ctx := context.Background()

	exporter, err := ConnectOTLP(ctx)
	if err != nil {
		panic(err)
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/", NewHandler(exporter))

	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		panic("could not start server")
	}
}
