package main

import (
	"net/http"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

// Declare directory where tasks are stored.
var taskPath = "tasks/"

// Declare valid "filepaths":
// - /plon/
// - /plon/styles.css
// - /plon/icon-192x192.png
var validFilePath = regexp.MustCompile(
	"^/plon/(styles.css|icon-192x192.png)?$",
)

// Declare valid paths to view/edit tasks:
// - /plon/view/<task id>
// - /plon/edit/<task id>
// - /plon/save/<task id>
var validIdPath = regexp.MustCompile(
	"^/plon/(view|edit|save)/([a-zA-Z0-9]+)$",
)

// Get valid filepaths from url.
// Return index handler for execution.
func MakeIndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validFilePath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			log.Error("Invalid path.")
			return
		}

		err := indexHandler(w, r, m[1]) // The file is the first subexpression.
		if err != nil {
			http.NotFound(w, r)
			log.Error(err)
			return
		}
	}
}

// Get valid task id(string) from url.
// Return apropriate handler for execution.
func MakeHandler(fn func(http.ResponseWriter, *http.Request, string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validIdPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			log.Error("Invalid page id.")
			return
		}

		err := fn(w, r, m[2]) // The id is the second subexpression.
		if err != nil {
			http.NotFound(w, r)
			log.Error(err)
			return
		}
	}
}

// Handle requests on /plon/.
// If not requesting a file, list all tasks from database by
// rendering the index.html template.
func indexHandler(w http.ResponseWriter, r *http.Request, fn string) error {
	if fn != "" {
		http.ServeFile(w, r, fn)
		log.WithFields(log.Fields{
			"file": fn,
		}).Info("Serving file.")
		return nil
	}

	p := &IndexPage{}
	for id, _ := range DB.Tasks {
		p.Alltasks = append(p.Alltasks,
			struct {
				Id    string
				Title string
			}{
				Id:    id,
				Title: DB.Tasks[id].Title,
			})
	}

	err := RenderTemplate(w, "index.html", p)
	if err != nil {
		return err
	}

	return nil
}

// Handle requests on /plon/add/.
// Load username list and render the add.html page.
func AddHandler(w http.ResponseWriter, r *http.Request) {
	p := &EditPage{
		Id: NewUID(), // Create new id.
	}

	for username, _ := range DB.Users {
		p.Addressees = append(p.Addressees, username)
	}

	err := RenderTemplate(w, "edit.html", p)
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}
}

// Handle requests on /plon/edit/.
// Load the task with specified id and render the edit.html template.
func EditHandler(w http.ResponseWriter, r *http.Request, id string) error {
	t := DB.Tasks[id]

	task, err := t.Read()
	if err != nil {
		return err
	}

	p := &EditPage{
		Id:         id,
		Title:      t.Title,
		Task:       task,
		Addressees: t.Addressees,
	}

	err = RenderTemplate(w, "edit.html", p)
	if err != nil {
		return err
	}

	return nil
}

// Handle requests on /plon/save/.
// Recieve task details through post method.
// Save the task.
// Redirect the user to view the task.
func SaveHandler(w http.ResponseWriter, r *http.Request, id string) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	t := Task{
		Title:      r.PostForm["title"][0],
		Path:       TaskPath + id + "/" + id + ".html",
		Addressees: r.PostForm["addressees"],
		Created:    time.Now(),
	}

	log.WithFields(log.Fields{
		"id":         id,
		"title":      t.Title,
		"addressees": t.Addressees,
	}).Info("Recieved new task.")

	err := t.Save(id, r.PostForm["task"][0])
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/plon/view/"+id, http.StatusFound)

	return nil
}

// Handle requests on /plon/view/.
// Read task with specified id from file and render task.html template.
func ViewHandler(w http.ResponseWriter, r *http.Request, id string) error {
	task, err := DB.Tasks[id].Read()
	if err != nil {
		return err
	}

	p := &ViewPage{
		Id:    id,
		Title: DB.Tasks[id].Title,
		Task:  string(task),
	}

	err = RenderTemplate(w, "task.html", p)
	if err != nil {
		return err
	}

	return nil
}
