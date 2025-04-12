package handler

import (
	"fmt"
	"net/http"
	"os"
	"vercel2otel/pkg/vercel"
	"vercel2otel/pkg/vercel2otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
)

// Handler is the Vercel entrypoint for this package.
func Handler(w http.ResponseWriter, r *http.Request) {
	auth := os.Getenv("HTTP_BASIC_SECRET")
	checksumSecret := os.Getenv("VERCEL_SECRET")
	format := os.Getenv("VERCEL_LINE_FORMAT")

	exporter, err := otlploghttp.New(r.Context())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not start new otlp exporter: %s", err.Error())
		w.WriteHeader(500)
		return
	}

	defer exporter.Shutdown(r.Context())

	handlerConfig := vercel2otel.Vercel2OtelHandlerConfig{
		Exporter:       exporter,
		Format:         vercel.ParserFormat(format),
		BasicAuth:      auth,
		ChecksumSecret: checksumSecret,
	}

	handlerConfig.HandleRequest(w, r)

	exporter.ForceFlush(r.Context())
}
