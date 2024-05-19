package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetMockEnv(key string) string {
	switch key {
	case "APP_ENV":
		return "test"
	case "APP_PORT":
		return "3000"
	case "APP_LOG_LEVEL":
		return "ERROR"
	default:
		return ""
	}
}

func TestHealthRoute(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	logBuf := new(bytes.Buffer)
	cfg := NewConfig(GetMockEnv, make([]string, 0))
	router := NewRouter(
		ctx,
		cfg,
		NewLogger(cfg, logBuf),
	)
	srv := httptest.NewServer(router)
	defer srv.Close()

	res, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatal("unable to complete Get request")
	}
	defer res.Body.Close()
	out, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err, "unable to read response data")
	}

	if logBuf.Len() > 0 {
		t.Fatalf("Error log should be empty, got: %d", logBuf.Len())
	}

	outStr := string(out)
	expected := "\"OK\"\n"
	if outStr != expected {
		t.Fatalf("Unexpected response from /health endpoint: %s (expected: %s)", outStr, expected)
	}
}
