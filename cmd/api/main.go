package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
)

type Application struct {
	logger      *slog.Logger
	templates   templateCache
	formDecoder *form.Decoder
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	templates, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	app := &Application{
		logger:      logger,
		templates:   templates,
		formDecoder: formDecoder,
	}

	srv := &http.Server{
		Addr:     ":8080",
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		Handler:  app.routes(),
	}

	logger.Info("starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
