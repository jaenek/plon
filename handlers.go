package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

// Declare directory where tasks are stored
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
var validIdPath = regexp.MustCompile(
	"^/plon/(view|edit|save)/([a-zA-Z0-9]+)$",
)

// Get valid filepaths from url.
func getFilename(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validFilePath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid path")
	}
	return m[1], nil // The file is the first subexpression.
}

// Get valid task id(string) from url.
func getId(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validIdPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Id")
	}
	return m[2], nil // The id is the second subexpression.
}

// Handle requests on /plon/.
// If not requesting a file, list all tasks from database by
// rendering the index.html template.
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

	err = RenderTemplate(w, "index.html", p)
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}
}

// Handle requests on /plon/view/.
// Read task with specified id from file and render task.html template.
func ViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		log.Error(err.Error())
		return
	}

	task, err := ioutil.ReadFile(DB.Tasks[id].Path)
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}

	p := &ViewPage{
		Id:    id,
		Title: DB.Tasks[id].Title,
		Task:  string(task),
	}

	err = RenderTemplate(w, "task.html", p)
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}
}

// Handle requests on /plon/add/.
// Load username list and render the add.html page.
func AddHandler(w http.ResponseWriter, r *http.Request) {
	p := &AddPage{
		Id: NewUID(), //
	}

	for username, _ := range DB.Users {
		p.Usernames = append(p.Usernames, username)
	}

	err := RenderTemplate(w, "add.html", p)
	if err != nil {
		http.NotFound(w, r)
		log.Error(err.Error())
		return
	}
}

// Handle requests on /plon/save/.
// Recieve task details through post method.
// Create task directory and save it to file.
// Add the task to the database and save it.
// Redirect the user to view the task.
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if err = r.ParseForm(); err != nil {
		log.Error(err.Error())
		http.NotFound(w, r)
		return
	}

	t := Task{
		Title:   r.PostForm["title"][0],
		Path:    taskPath + id + "/" + id + ".html",
		Created: time.Now(),
	}

	path := taskPath + id
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			log.Error(err.Error())
			http.NotFound(w, r)
			return
		}
	}

	task := []byte(r.PostForm["task"][0])
	err = ioutil.WriteFile(t.Path, task, 0644)
	if err != nil {
		log.Error(err.Error())
		http.NotFound(w, r)
		return
	}

	for _, username := range r.PostForm["addressees"] {
		user := DB.Users[username]
		user.Tasks = append(DB.Users[username].Tasks, id)
		DB.Users[username] = user
	}
	DB.Tasks[id] = t
	DB.Write()

	log.WithFields(log.Fields{
		"id":         id,
		"title":      t.Title,
		"addressees": r.PostForm["addressees"],
	}).Info("Recieved new task.")

	http.Redirect(w, r, "/plon/view/"+id, http.StatusFound)
}
