package main

import (
	"errors"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
)

var validFilePath = regexp.MustCompile(
	"^/plon/(styles.css|icon-192x192.png)?$",
)

var validIdPath = regexp.MustCompile(
	"^/plon/(view|edit)/([a-zA-Z0-9]+)$",
)

func getFilename(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validFilePath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid path")
	}
	return m[1], nil // The file is the first subexpression.
}

func getId(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validIdPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Id")
	}
	return m[2], nil // The id is the second subexpression.
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fn, err := getFilename(w, r)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if fn != "" {
		http.ServeFile(w, r, fn)
		log.WithFields(log.Fields{
			"file": fn,
		}).Info("Serving file.")
		return
	}

	db, err := ReadDatabase("db.json")
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}

	p := &IndexPage{}
	for id, _ := range db.Tasks {
		p.Alltasks = append(p.Alltasks,
			struct {
				Id    string
				Title string
			}{
				Id:    id,
				Title: db.Tasks[id].Title,
			})
	}

	err = RenderTemplate(w, "index.html", p)
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		log.Error(err.Error())
		return
	}

	db, err := ReadDatabase("db.json")
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}

	p := &ViewPage{
		Id:    id,
		Title: db.Tasks[id].Title,
		Task: "This is a test task later it would be loaded from file (" +
			db.Tasks[id].Path +
			")rendered from markdown to html",
	}

	err = RenderTemplate(w, "task.html", p)
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}
}
