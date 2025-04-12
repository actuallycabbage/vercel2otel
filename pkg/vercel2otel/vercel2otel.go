package vercel2otel

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"vercel2otel/pkg/formatter"
	"vercel2otel/pkg/vercel"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
)

type Vercel2OtelHandlerConfig struct {
	Format         vercel.ParserFormat
	Exporter       *otlploghttp.Exporter
	BasicAuth      string
	ChecksumSecret string
}

func (v Vercel2OtelHandlerConfig) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if !basicAuthCheck(w, r, v.BasicAuth) {
		return
	}

	// expected checksum
	checksum := r.Header.Get("x-vercel-signature")
	if checksum == "" {
		fmt.Fprintf(os.Stderr, "missing vercel signature")
		w.WriteHeader(400)
		return
	}

	// vercel computes these via sha1
	hash := hmac.New(sha1.New, []byte(v.ChecksumSecret))

	// tee into the line reader and hmac
	defer r.Body.Close()
	reader := io.TeeReader(r.Body, hash)

	// prep parser
	parser, err := vercel.GetParser(v.Format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected format: %s\n", err.Error())
		w.WriteHeader(500)
		return
	}

	// parse
	lines, err := parser(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parser error: %s\n", err.Error())
		w.WriteHeader(400)
		return
	}

	// // drain the reader
	// io.ReadAll(reader)

	// verify sum
	if hex.EncodeToString(hash.Sum(nil)) != checksum {
		fmt.Fprintf(os.Stderr, "bad checksum\n")
		w.WriteHeader(400)
		return
	}

	// send of the lines
	err = v.Exporter.Export(r.Context(), formatter.TransformLogItems(lines))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not export lines: %s\n", err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func basicAuthCheck(w http.ResponseWriter, r *http.Request, token string) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(400)
		return false
	}
	userAuth, found := strings.CutPrefix(authHeader, "Basic ")
	if !found {
		w.WriteHeader(400)
		return false
	}

	if !strings.EqualFold(token, userAuth) {
		if !found {
			w.WriteHeader(401)
			return false
		}
	}

	return true
}
