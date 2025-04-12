package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"vercel2otel/pkg/vercel"
	"vercel2otel/pkg/vercel2otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
)

func main() {
	auth := os.Getenv("HTTP_BASIC_SECRET")
	checksumSecret := os.Getenv("VERCEL_SECRET")
	format := os.Getenv("VERCEL_LINE_FORMAT")

	ctx := context.Background()

	exporter, err := otlploghttp.New(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not start new otlp exporter: %s", err.Error())
		os.Exit(1)
	}

	defer exporter.Shutdown(ctx)

	handlerConfig := vercel2otel.Vercel2OtelHandlerConfig{
		Exporter:       exporter,
		Format:         vercel.ParserFormat(format),
		BasicAuth:      auth,
		ChecksumSecret: checksumSecret,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlerConfig.HandleRequest)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic("could not start server")
	}
}
