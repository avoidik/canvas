package server

import (
	"canvas/handlers"
	"canvas/model"
	"context"
)

type SignUpMock struct{}

func (s *Server) setupRoutes() {
	handlers.Health(s.mux)
	handlers.FrontPage(s.mux)
	handlers.NewsletterSignup(s.mux, &SignUpMock{})
	handlers.NewsletterThanks(s.mux)
}

func (s *SignUpMock) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	return "", nil
}
