package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
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

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
