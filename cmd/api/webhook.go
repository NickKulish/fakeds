package main

import (
	//"encoding/json"
	//"io"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"fakeds.nkulish.dev/ui"
)

type Document struct {
	ID       string
	Name     string
	PDFBytes string
}

type Recipient struct {
	ID     string
	Signed string
}

type WebhookData struct {
	Recipients []Recipient
	EnvelopeID string
	Documents  []Document
}

func ReadFakePDF() (string, error) {
	bytes, err := ui.Files.ReadFile("pdf/fake.pdf")
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func NewWebhookData(data EnvelopeData) (*WebhookData, error) {
	now := time.Now().Format("2006-01-02T15:04:05.000")
	recipients := []Recipient{}
	for _, p := range data.Participants {
		recipients = append(recipients, Recipient{
			ID:     p.ID,
			Signed: now,
		})
	}
	documents := []Document{}
	fakePDF, err := ReadFakePDF()
	if err != nil {
		return nil, err
	}
	for _, d := range data.Workflow.Attachments {
		if d.DocusignID != "" {
			documents = append(documents, Document{
				ID:       d.DocusignID,
				Name:     "fake.pdf",
				PDFBytes: fakePDF,
			})
		}
	}
	return &WebhookData{
		Recipients: recipients,
		EnvelopeID: data.Envelope.EnvelopeID,
		Documents:  documents,
	}, nil
}

func (app *Application) sendWebhook(webhookUrl string, data EnvelopeData) (int, error) {
	wd, err := NewWebhookData(data)
	if err != nil {
		return 0, err
	}

	buf, err := app.renderXML(*wd)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(
		webhookUrl,
		"text/xml; charset=utf-8",
		buf)

	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

type webhookCreateForm struct {
	Validator   `form:"-"`
	EnvelopeUrl string `form:"envelopeUrl"`
	WebhookUrl  string `form:"webhookUrl"`
	Token       string `form:"token"`
	EnvelopeID  string `form:"envelopeID"`
}

func (app *Application) webhookCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = webhookCreateForm{
		EnvelopeUrl: "http://192.168.1.101:8082/workflow/c/nkulish/docusign-envelope-info/",
		WebhookUrl:  "https://2394-119-92-53-249.ngrok-free.app/docusign/c/nkulish/webhook",
	}

	app.render(w, r, http.StatusOK, "webhook.tmpl", data)
}

func (app *Application) webhookCreatePost(w http.ResponseWriter, r *http.Request) {
	var form webhookCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data, err := app.getEnvelopeData(form.EnvelopeUrl, form.EnvelopeID, form.Token)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	code, err := app.sendWebhook(form.WebhookUrl, *data)

	w.Write([]byte(strconv.Itoa(code)))
}
