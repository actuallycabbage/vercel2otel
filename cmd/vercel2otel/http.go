package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"vercel2otel/pkg/vercel"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
)

func NewHandler(exporter *otlploghttp.Exporter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("could not copy body: %w", err)
		}

		// // check the sha
		// ok, err := vercel.VerifyBody(bytes.NewReader(b), []byte{}, "")
		// if err != nil {
		// 	log.Fatalf("could not verify body: %w", err)
		// }

		// if !ok {
		// 	w.WriteHeader(401)
		// 	return
		// }

		// parse the body
		logs, err := vercel.ParseJSON(bytes.NewReader(b))
		if err != nil {
			log.Fatalf("could not parse body: %w", err)
		}

		// emit
		exporter.Export(r.Context(), TransformLogItems(logs))

		fmt.Println(logs)

		w.WriteHeader(200)
	}
}
