package main

import (
	"net/http"
	"text/template"

	log "github.com/sirupsen/logrus"
)

// Declare directory where template files are stored.
const TemplatePath = "templates/"

type UserPage struct {
	Username string
	Alltasks []struct {
		Id    string
		Title string
	}
}

type EditPage struct {
	Id         string
	Title      string
	Task       string
	Usernames  []string
	Addressees map[string]bool
	Due        string
	Deletable  bool
}

type ViewPage struct {
	Id    string
	Title string
	Task  string
}

// Render template from filepath using passed data.
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
