package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
)

type Task struct {
	Id      int
	Path    string
	Created time.Time
	Due     time.Time
}

type User struct {
	Tasks map[string]Task
}

type Database struct {
	Users map[string]User
}

type Page struct {
	Alltasks []string
}

func (db *Database) write() error {
	b, err := json.MarshalIndent(db, "", "\t")
	if err != nil {
		return err
	}

	path := "db.json"
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"file": path,
	}).Info("Saving to file.")
	return nil
}

func read(path string) (*Database, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"file": path,
	}).Info("Unmarshalling json.")

	db := &Database{}
	err = json.Unmarshal(b, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) error {
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/plon"):]
	if path != "/" {
		http.ServeFile(w, r, path)
		log.WithFields(log.Fields{
			"file": path,
		}).Info("Serving file.")
		return
	}

	db, err := read("db.json")
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

	err = renderTemplate(w, "index.html", p)
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}
}

func main() {
	http.HandleFunc("/plon/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
