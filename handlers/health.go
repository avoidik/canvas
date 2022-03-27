package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

func Health(mux chi.Router, p Pinger) {
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := p.Ping(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	})
}
