package main

import (
	"embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"undevised.com/recipes/templates"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	mux := http.NewServeMux()
	mux.Handle("GET /", logRequest(logger, handleIndex(templates.FS)))

	logger.Info("Starting web server")

	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func logRequest(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recordedWriter := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(recordedWriter, r)

		logger.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", recordedWriter.status),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

func handleIndex(fs embed.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(fs, "index.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	})
}
