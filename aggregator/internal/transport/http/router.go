package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/qsoulior/news/aggregator/internal/service"
	"github.com/qsoulior/news/aggregator/internal/transport/http/handler"
	"github.com/rs/cors"
)

func NewRouter(service service.News) http.Handler {
	mux := chi.NewMux()
	mux.Use(cors.Default().Handler)
	mux.Use(middleware.AllowContentType("application/json"))
	mux.Use(middleware.RealIP)
	mux.Use(LoggerMiddleware())
	mux.Use(RecovererMiddleware())

	news := handler.NewNews(service)
	mux.Get("/news", news.List)
	mux.Get("/news/{id}", news.Get)

	return mux
}
