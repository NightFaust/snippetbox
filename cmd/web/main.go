package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", ".ui/static", "Path to static assets")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", cfg.addr))

	err := http.ListenAndServe(cfg.addr, app.routes(cfg.staticDir))
	logger.Error(err.Error())
	os.Exit(1)
}

type config struct {
	addr      string
	staticDir string
}

type application struct {
	logger *slog.Logger
}
