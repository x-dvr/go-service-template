package json_test

import (
	"bytes"
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

type testVal struct {
	name     string
	val      any
	expected []byte
}

func TestJsonEncode(t *testing.T) {
	values := make([]testVal, 0, 2)
	values = append(values,
		testVal{name: "encode string", val: "some string", expected: []byte("\"some string\"\n")},
		testVal{name: "encode int", val: 42, expected: []byte("42\n")},
	)

	for _, test := range values {
		t.Run(test.name, func(t *testing.T) {
			// t.Parallel()
			mw := newMockWriter()

			err := json.Encode(mw, 123, test.val)
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

			if !bytes.Equal(mw.buf, test.expected) {
				t.Fatalf("Should write: %v , written: %v", test.expected, mw.buf)
			}
		})
	}
}
