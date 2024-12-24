package main

import (
	"net/http"

	"fakeds.nkulish.dev/ui"
	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.FileServerFS(ui.Files)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/", app.webhookCreate)
	router.HandlerFunc(http.MethodPost, "/", app.webhookCreatePost)

	return router
}
