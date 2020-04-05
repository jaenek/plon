package main

import (
	"net/http"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type Page struct {
	Name     string
	Alltasks []string
}

func RenderTemplate(w http.ResponseWriter, tmpl string, p interface{}) error {
	log.WithFields(log.Fields{
		"file": tmpl,
	}).Info("Rendering template.")

	t, err := template.ParseFiles(tmpl)
	if err != nil {
		return err
	}

	err = t.Execute(w, p)
	if err != nil {
		return err
	}

	return nil
}
