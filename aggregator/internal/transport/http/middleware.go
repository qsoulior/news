package http

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Middleware func(next http.Handler) http.Handler

type loggerWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *loggerWriter) Write(bytes []byte) (int, error) {
	n, err := w.ResponseWriter.Write(bytes)
	w.size += n
	return n, err
}

func (w *loggerWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func LoggerMiddleware(logger *zerolog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &loggerWriter{w, 200, 0}
			start := time.Now()
			next.ServeHTTP(writer, r)
			logger.Info().
				Str("addr", r.RemoteAddr).
				Str("method", r.Method).
				Stringer("url", r.URL).
				Str("proto", r.Proto).
				Int("status", writer.status).
				Int("size", writer.size).
				Dur("duration", time.Since(start)).
				Send()
		})
	}
}

func RecovererMiddleware(logger *zerolog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					if r == http.ErrAbortHandler {
						panic(r)
					}

					w.WriteHeader(http.StatusInternalServerError)

					err, ok := r.(error)
					if ok {
						logger.Error().Err(err).Send()
						return
					}

					logger.Error().Msgf("%s", r)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
