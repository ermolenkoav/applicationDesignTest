package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"applicationDesignTest/internal/logg"
	"applicationDesignTest/internal/model"
)

type bookingService interface {
	DoBookingOrder(context.Context, model.Order) error
}

type api struct {
	bService bookingService
	router   *chi.Mux
	srv      *http.Server
}

func NewServe(bService bookingService) *api {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	a := &api{
		bService: bService,
		router:   mux}

	mux.Post("/orders", a.createOrder)
	return a
}

func (a *api) ListenAndServe() error {
	a.srv = &http.Server{
		Addr:         ":8080",
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
	logg.Info("server listening on " + a.srv.Addr)
	return http.ListenAndServe(":8080", a.router)
}

func (a *api) Shutdown(ctx context.Context) error {
	logg.Info("API-Shutdown")
	return a.srv.Shutdown(ctx)
}
