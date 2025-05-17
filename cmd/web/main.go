package main

import (
	"embed"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"undevised.com/recipes/ui"
)

type application struct {
	logger    *slog.Logger
	templates embed.FS
}

func main() {
	addr := flag.String("addr", ":3000", "TCP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	app := &application{
		logger:    logger,
		templates: ui.Files,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", http.NotFound)
	mux.HandleFunc("GET /{$}", app.home)

	logger.Info("Starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.logRequest(mux))
	if err != nil {
		log.Fatal(err)
	}
}
