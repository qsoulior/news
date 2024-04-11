package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/qsoulior/news/aggregator/internal/service"
	"github.com/qsoulior/news/aggregator/internal/transport/http/handler"
	"github.com/rs/zerolog"
)

func NewRouter(logger *zerolog.Logger, service service.News) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.AllowContentType("application/json"))
	mux.Use(middleware.RealIP)
	mux.Use(LoggerMiddleware(logger))
	mux.Use(RecovererMiddleware(logger))

	news := handler.NewNews(handler.NewsConfig{
		Logger:  logger,
		Service: service,
	})
	mux.Get("/news", news.Get)

	return mux
}
