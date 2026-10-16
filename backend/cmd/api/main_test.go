package main

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthRoutes(t *testing.T) {
	for _, test := range []struct {
		name       string
		method     string
		path       string
		pingError  error
		wantStatus int
		wantBody   string
	}{
		{"liveness without database", "GET", "/healthz", errors.New("private connection URL"), 200, `{"status":"ok"}`},
		{"ready database", "GET", "/readyz", nil, 200, `{"status":"ready"}`},
		{"unavailable database", "GET", "/readyz", errors.New("private connection URL"), 503, `{"status":"unavailable"}`},
		{"unknown route", "GET", "/sessions", nil, 404, "404 page not found"},
		{"unsupported method", "POST", "/healthz", nil, 405, "Method Not Allowed"},
	} {
		t.Run(test.name, func(t *testing.T) {
			handler := newHandler(func(context.Context) error { return test.pingError })
			response, body := testRequest(t, handler, test.method, test.path)
			if response.StatusCode != test.wantStatus || strings.TrimSpace(string(body)) != test.wantBody {
				t.Fatalf("got %d %q", response.StatusCode, string(body))
			}
			if strings.Contains(string(body), "private") {
				t.Fatal("database details exposed")
			}
		})
	}
}

func TestServerStopsOnCancellation(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- serve(ctx, listener, newHandler(func(context.Context) error { return nil })) }()
	response, err := http.Get("http://" + listener.Addr().String() + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("server did not stop")
	}
}

func testRequest(t *testing.T, app *fiber.App, method, path string) (*http.Response, []byte) {
	t.Helper()
	response, err := app.Test(httptest.NewRequest(method, path, nil))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	return response, body
}
