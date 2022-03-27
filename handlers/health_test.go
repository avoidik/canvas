package handlers_test

import (
	"canvas/handlers"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"
)

type PingerMock struct {
	err error
}

func (p *PingerMock) Ping(ctx context.Context) error {
	return p.err
}

func TestHealth(t *testing.T) {
	t.Run("return 200", func(t *testing.T) {
		is := is.New(t)

		mux := chi.NewMux()
		handlers.Health(mux, &PingerMock{})
		code, _, _ := makeGetRequest(mux, "/health")
		is.Equal(http.StatusOK, code)
	})

	t.Run("return 502", func(t *testing.T) {
		is := is.New(t)

		mux := chi.NewMux()
		handlers.Health(mux, &PingerMock{err: errors.New("error")})
		code, _, _ := makeGetRequest(mux, "/health")
		is.Equal(http.StatusBadGateway, code)
	})
}

func makeGetRequest(handler http.Handler, target string) (int, http.Header, string) {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	result := res.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, result.Header, string(bodyBytes)
}
