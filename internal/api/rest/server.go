package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "applicationDesignTest/docs"
	"applicationDesignTest/internal/logg"
	"applicationDesignTest/internal/model"
)

const addr = ":8080"

type bookingService interface {
	DoBookingOrder(context.Context, model.Order) error
}

type Server struct {
	srv      *http.Server
	bService bookingService
}

//	@title			Booking API
//	@version		1.0
//	@description	Hotel room booking service.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Alexey Ermolenko
//	@contact.email	ermolenkoav@gmail.com

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func NewServer(bService bookingService) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	s := &Server{bService: bService}

	r.Post("/api/v1/orders", s.createOrder)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"),
	))

	s.srv = &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return s
}

func (s *Server) ListenAndServe() error {
	logg.Info("server listening on %s", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	logg.Info("shutting down")
	return s.srv.Shutdown(ctx)
}
