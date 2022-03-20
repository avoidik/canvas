package handlers_test

import (
	"canvas/handlers"
	"canvas/model"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"
)

type SignUpMock struct {
	email model.Email
}

func (s *SignUpMock) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	s.email = email
	return "", nil
}

func TestNewsletterSignup(t *testing.T) {
	mux := chi.NewMux()
	s := &SignUpMock{}
	handlers.NewsletterSignup(mux, s)

	t.Run("valid email address", func(t *testing.T) {
		is := is.New(t)
		code, _, _ := makePostRequest(t, mux, "/newsletter/signup", createFormHeader(), strings.NewReader("email=me%40example.com"))
		is.Equal(http.StatusFound, code)
		is.Equal(model.Email("me@example.com"), s.email)
	})

	t.Run("invalid email address", func(t *testing.T) {
		is := is.New(t)
		code, _, _ := makePostRequest(t, mux, "/newsletter/signup", createFormHeader(), strings.NewReader("email=notanemail"))
		is.Equal(http.StatusBadRequest, code)
	})
}

func makePostRequest(t *testing.T, handler http.Handler, target string, header http.Header, body io.Reader) (int, http.Header, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, target, body)
	req.Header = header
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	result := res.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	return result.StatusCode, result.Header, string(bodyBytes)
}

func createFormHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	return header
}
