package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"applicationDesignTest/internal/logg"
	"applicationDesignTest/internal/model"
)

type repo interface {
	GetAvailability() ([]model.RoomAvailability, error)
	SaveOrder(model.Order) error
}

type api struct {
	router *chi.Mux
	srv    *http.Server
	repo   repo
}

func NewServe(repo repo) *api {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	a := &api{router: mux,
		repo: repo}

	mux.Post("/orders", a.createOrder)

	return a
}

func (a *api) ListenAndServe() error {
	a.srv = &http.Server{
		Addr: ":8080",
	}
	logg.Info("server listening on " + a.srv.Addr)
	return http.ListenAndServe(":8080", a.router)
}

func (a *api) Shutdown(ctx context.Context) error {
	logg.Info("API-Shutdown")
	return a.srv.Shutdown(ctx)
}
