package services

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/BX-Research/ngadu-monorepo/backend/app/types"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type finderStub struct {
	err error
	id  string
}

func (s *finderStub) FindUser(_ context.Context, id string) (types.UserInfo, error) {
	s.id = id
	return types.UserInfo{ID: id, DisplayName: "Blake Demo"}, s.err
}

func TestFindUser(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"found", nil, 200},
		{"missing", pgx.ErrNoRows, 404},
		{"database failure", errors.New("database unavailable"), 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stub := &finderStub{err: tc.err}
			app := fiber.New()
			app.Get("/users/:id", FindUser(stub))
			res, err := app.Test(httptest.NewRequest("GET", "/users/00000000-0000-4000-8000-000000000002", nil))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.status {
				t.Fatalf("status = %d, want %d", res.StatusCode, tc.status)
			}
			if stub.id != "00000000-0000-4000-8000-000000000002" {
				t.Fatalf("lookup received ID %q", stub.id)
			}
		})
	}
}
