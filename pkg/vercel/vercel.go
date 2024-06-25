package vercel

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

type LogLevel string

const (
	StdErr LogLevel = "stderr"
	StdOut LogLevel = "stdout"
)

type LogSource string

const (
	BuildLogSource    LogSource = "build"
	StaticLogSource   LogSource = "static"
	ExternalLogSource LogSource = "external"
	LambdaLogSource   LogSource = "lambda"
)

type LogItem struct {
	ID      string `json:"id"`
	Message string `json:"message"`

	// TODO: this might actually be an int
	Timestamp    int64     `json:"timestamp"`
	Type         string    `json:"type,omitempty"`
	Source       LogSource `json:"source"`
	ProjectID    string    `json:"projectId"`
	DeploymentID string    `json:"deploymentId"`
	BuildID      string    `json:"buildId"`
	Host         string    `json:"host"`
	Entrypoint   string    `json:"entrypoint,omitempty"`
	RequestID    string    `json:"requestId,omitempty"`

	// this might be a signed int
	StatusCode      int      `json:"statusCodee,omitempty"`
	Destination     string   `json:"destination,omitempty"`
	Path            string   `json:"path,omitempty"`
	ExecutionRegion string   `json:"executionRegion,omitempty"`
	Level           LogLevel `json:"level,omitempty"`
	Proxy           struct {
		// again, could be an int.
		Timestamp int64  `json:"timestamp"`
		Method    string `json:"method"`
		Scheme    string `json:"scheme"`
		Host      string `json:"host"`
		Path      string `json:"path"`
		UserAgent string `json:"userAgent"`
		Referer   string `json:"referer"`
		// again, could be an int
		StatusCode  int    `json:"statusCode"`
		ClientIP    string `json:"clientIp"`
		Region      string `json:"region"`
		CacheID     string `json:"cacheId"`
		VercelCache string `json:"vercelCache"`
	} `json:"proxy,omitempty"`
}

type LogParser interface {
	Read(io.Reader) ([]LogItem, error)
}

func ParseJSON(reader io.Reader) (logs []LogItem, err error) {
	err = json.NewDecoder(reader).Decode(&logs)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("could not decode json: %w", err)
	}

	return logs, nil
}

func ParseNDJSON(reader io.Reader) (logs []LogItem, err error) {
	decoder := json.NewDecoder(reader)

	// TODO: not sure if More() is the right one to use here.
	for decoder.More() {
		var item LogItem
		err = decoder.Decode(&item)
		if err != nil {
			return nil, fmt.Errorf("could not parse ndjson item: %w", err)
		}

		logs = append(logs, item)
	}

	return logs, nil
}

func VerifyBody(reader io.Reader, secret []byte, expectedHash string) (bool, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return false, fmt.Errorf("could not read body: %w", err)
	}

	h := hmac.New(sha1.New, secret)
	_, err = h.Write(body)
	if err != nil {
		return false, fmt.Errorf("could not write hmac: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)) == expectedHash, nil
}
