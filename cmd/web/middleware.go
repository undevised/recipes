package main

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recordedWriter := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(recordedWriter, r)

		app.logger.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("protocol", r.Proto),
			slog.String("client", r.RemoteAddr),
			slog.Int("status", recordedWriter.status),
			slog.Duration("duration", time.Since(start)),
		)
	})
}
