package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
)

var validPath = regexp.MustCompile(
	"^/plon/(((view|edit)/([a-zA-Z0-9]+))|styles.css|icon-192x192.png)?$",
)

func getName(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Name")
	}
	return m[2], nil // The name is the second subexpression.
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	name, err := getName(w, r)
	if err != nil {
		log.Error(err.Error())
		return
	}

	fmt.Printf("\"%s\"\n", name)
	if name != "" {
		http.ServeFile(w, r, name)
		log.WithFields(log.Fields{
			"file": name,
		}).Info("Serving file.")
		return
	}

	db, err := ReadDatabase("db.json")
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}

	p := &Page{
		Alltasks: []string{},
	}
	for _, user := range db.Users {
		for title, _ := range user.Tasks {
			p.Alltasks = append(p.Alltasks, title)
		}
	}

	err = RenderTemplate(w, "index.html", p)
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}
}
