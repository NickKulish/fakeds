package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	tt "text/template"

	"fakeds.nkulish.dev/ui"
)

type templateData struct {
	Form   any
	Result string
}

func (app *Application) newTemplateData() templateData {
	return templateData{}
}

type templateCache map[string]*template.Template

func newTemplateCache() (templateCache, error) {
	cache := templateCache{}

	// html
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.tmpl",
			page,
		}
		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templates[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.logger.Error(err.Error())
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *Application) renderXML(data WebhookData) (*bytes.Buffer, error) {
	patterns := []string{
		"xml/webhook.xml",
	}
	ts, err := tt.ParseFS(ui.Files, patterns...)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = ts.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
