package environmentheader_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nilskohrs/environmentheader"
)

func TestAddingRequestHeader(t *testing.T) {
	value := "FOO BAR"
	header := "Test-Header"
	envVar := "TEST_ENV"

	t.Setenv(envVar, value)

	cfg := environmentheader.CreateConfig()
	cfg.RequestHeaders = append(cfg.RequestHeaders, environmentheader.HeaderMapping{
		Header: header,
		Env:    envVar,
	})

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := environmentheader.New(ctx, next, cfg, "environmentheader")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://test", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)

	assertRequestHeaderAdded(t, req, header, value)
}

func TestAddingResponseHeader(t *testing.T) {
	value := "FOO"
	header := "Test-Header"
	envVar := "TEST_ENV"

	err := os.Setenv(envVar, value)
	if err != nil {
		t.Fatal(err)
	}

	cfg := environmentheader.CreateConfig()
	cfg.ResponseHeaders = append(cfg.ResponseHeaders, environmentheader.HeaderMapping{
		Header: header,
		Env:    envVar,
	})

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := environmentheader.New(ctx, next, cfg, "environmentheader")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://test", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)

	assertResponseHeaderAdded(t, recorder, header, value)
}

func assertRequestHeaderAdded(t *testing.T, req *http.Request, header, value string) {
	t.Helper()
	if req.Header.Get(header) != value {
		t.Errorf("Expected header '%s' to be '%s', but was '%s'", header, value, req.Header.Get(header))
	}
}

func assertResponseHeaderAdded(t *testing.T, resp *httptest.ResponseRecorder, header, value string) {
	t.Helper()
	if resp.Header().Get(header) != value {
		t.Errorf("Expected header '%s' to be '%s', but was '%s'", header, value, resp.Header().Get(header))
	}
}
