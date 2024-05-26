package json_test

import (
	"net/http"
	"testing"

	"github.com/x-dvr/go-service-template/cmd/api/internal/json"
)

type mockWriter struct {
	h      http.Header
	buf    []byte
	status int
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		h:   make(http.Header, 1),
		buf: make([]byte, 0, 100),
	}
}

func (mw *mockWriter) Header() http.Header {
	return mw.h
}

func (mw *mockWriter) Write(bytes []byte) (int, error) {
	mw.buf = append(mw.buf, bytes...)
	return len(bytes), nil
}

func (mw *mockWriter) WriteHeader(statusCode int) {
	mw.status = statusCode
}

func TestJsonEncodeString(t *testing.T) {
	mw := newMockWriter()
	val := "some string"
	err := json.Encode(mw, 123, val)
	if err != nil {
		t.Fatalf("Should not return an error, got: %v", err)
	}

	if mw.status != 123 {
		t.Fatalf("Should write status 123, got: %v", mw.status)
	}

	if len(mw.h) != 1 {
		t.Fatalf("Should write 1 header, got: %v", len(mw.h))
	}

	if mw.h.Get("Content-Type") != "application/json" {
		t.Fatalf("Should write content-type header as 'application/json', got: %v", mw.h.Get("Content-Type"))
	}

	if string(mw.buf) != "\"some string\"\n" {
		t.Fatalf("Should write string: %v , written: %v", []byte("\"some string\""), mw.buf)
	}
}
