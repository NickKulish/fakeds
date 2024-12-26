package main

import (
	"errors"
	"net/http"
)

type Envelope struct {
	EnvelopeID string `json:"docusign-envelope/envelope-id"`
}

type Attachment struct {
	ID         string `json:"attachment-id"`
	DocumentID string `json:"document-id"`
	DocusignID string `json:"docusign-document-id"`
}

type Workflow struct {
	ID          string       `json:"id"`
	Attachments []Attachment `json:"attachments"`
}

type Participant struct {
	ID string `json:"participant-id"`
}

type EnvelopeData struct {
	Envelope     Envelope      `json:"docusign-envelope"`
	Workflow     Workflow      `json:"workflow"`
	Participants []Participant `json:"participants"`
	Error        string        `json:"error"`
}

func (app *Application) getEnvelopeData(url, envelopeID, token string) (*EnvelopeData, error) {
	req, err := http.NewRequest(http.MethodGet, url+envelopeID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data EnvelopeData
	err = app.readJSON(resp.Body, &data)
	app.logger.Info("ENVELOPE DATA", "data", data)
	if err != nil {
		return nil, err
	}
	if data.Error != "" {
		return nil, errors.New(data.Error)
	}

	return &data, nil
}
