package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"vercel2otel/pkg/formatter"
	"vercel2otel/pkg/vercel"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	// check auth header
	auth := os.Getenv("HTTP_BASIC_SECRET")
	if auth != "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(400)
			return
		}
		userAuth, found := strings.CutPrefix(authHeader, "Basic ")
		if !found {
			w.WriteHeader(400)
			return
		}

		if !strings.EqualFold(auth, userAuth) {
			if !found {
				w.WriteHeader(401)
				return
			}
		}
	}

	exporter, err := otlploghttp.New(r.Context())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not start new otlp exporter: %s", err.Error())
		w.WriteHeader(500)
		return
	}
	defer exporter.Shutdown(r.Context())
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not copy body: %s", err.Error())
		w.WriteHeader(500)
		return
	}

	logs, err := vercel.ParseJSON(bytes.NewReader(b))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse body: %s\n", err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = exporter.Export(r.Context(), formatter.TransformLogItems(logs))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse body: %s\n", err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)

}
