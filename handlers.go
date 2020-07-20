package main

import (
	"net/http"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

// Declare public directory path:
const PublicPath = "public/"

// Declare valid "filepaths":
// - /plon
// - /plon/
// - /plon/styles.css
// - /plon/markdown.css
// - /plon/plon.webmanifest
// - /plon/icons/favicon.png
// - /plon/icons/icon-<w>x<h>.png
// - /plon/fonts/<fontname>.ttf
var validFilePath = regexp.MustCompile(
	"^/plon(/|(/styles|/markdown|/light|/dark).css|/cookies.js|/plon.webmanifest|/icons/(favicon|icon-[0-9]+x[0-9]+).png|/fonts/[a-zA-Z]+-(Bold|Regular).ttf)?$",
)

// Declare valid paths to view/edit/save/delete tasks:
// - /plon/view/<task id>
// - /plon/edit/<task id>
// - /plon/save/<task id>
var validIdPath = regexp.MustCompile(
	"^/plon/(view|edit|save|delete)/([a-zA-Z0-9]+)$",
)

// Declare valid user paths:
var validUserPath = regexp.MustCompile(
	"^/plon/(user|profile)/([a-zA-Zęóąśłżźćń]+)$",
)

// Get valid filepaths from url.
// Return index handler for execution.
func MakeIndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validFilePath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			log.WithFields(log.Fields{
				"path": r.URL.Path,
			}).Error("Invalid path.")
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
func MakeIdHandler(fn func(http.ResponseWriter, *http.Request, string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validIdPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			log.WithFields(log.Fields{
				"path": r.URL.Path,
			}).Error("Invalid page id.")
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

// Get valid username from url.
// Return apropriate handler for execution.
func MakeUserHandler(fn func(http.ResponseWriter, *http.Request, []string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validUserPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			log.WithFields(log.Fields{
				"path": r.URL.Path,
			}).Error("Invalid username.")
			return
		}

		err := fn(w, r, m[2:]) // The username is the second subexpression.
		if err != nil {
			http.NotFound(w, r)
			log.Error(err)
			return
		}
	}
}

// Handle requests on /plon/user/.
// If not requesting a file redirect to /plon/user/wszyscy
func indexHandler(w http.ResponseWriter, r *http.Request, fn string) error {
	if fn != "" && fn != "/" {
		http.ServeFile(w, r, PublicPath+fn)
		log.WithFields(log.Fields{
			"file": fn,
		}).Info("Serving file.")
		return nil
	}

	http.Redirect(w, r, "/plon/user/wszyscy", http.StatusFound)

	return nil
}

// Handle requests on /plon/add/.
// Load username list and render the add.html page.
func AddHandler(w http.ResponseWriter, r *http.Request) {
	p := &EditPage{
		Id:        NewUID(), // Create new id.
		Deletable: false,
	}

	for username, _ := range DB.Users {
		p.Usernames = append(p.Usernames, username)
	}

	err := RenderTemplate(w, TemplatePath+"edit.html", p)
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

	task, err := t.Read(".md")
	if err != nil {
		return err
	}

	p := &EditPage{
		Id:         id,
		Title:      t.Title,
		Task:       task,
		Addressees: t.Addressees,
		Deletable:  true,
	}

	for username, _ := range DB.Users {
		if DB.Users[username].Tasks[id] != true {
			p.Usernames = append(p.Usernames, username)
		}
	}

	err = RenderTemplate(w, TemplatePath+"edit.html", p)
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

	addressees := map[string]bool{}
	if len(r.PostForm["addressees"]) == 0 {
		addressees["wszyscy"] = true
	} else {
		for _, username := range r.PostForm["addressees"] {
			addressees[username] = true
		}
	}

	t := Task{
		Title:      r.PostForm["title"][0],
		Path:       TaskPath + id + "/" + id,
		Addressees: addressees,
		Created:    time.Now(),
	}

	log.WithFields(log.Fields{
		"id":         id,
		"title":      t.Title,
		"addressees": r.PostForm["addressees"],
	}).Info("Recieved new task.")

	err := t.Save(id, r.PostForm["task"][0])
	if err != nil {
		return err
	}

	DB.Tasks[id] = t
	err = WriteJSON("db.json", DB)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/plon/view/"+id, http.StatusFound)

	return nil
}

// Handle requests on /plon/delete/.
// Delte  the task with specified id and redirect to index.html.
func DeleteHandler(w http.ResponseWriter, r *http.Request, id string) error {
	err := DB.Tasks[id].Delete(id)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/plon/", http.StatusFound)

	return nil
}

// Handle requests on /plon/view/.
// Read task with specified id from file and render task.html template.
func ViewHandler(w http.ResponseWriter, r *http.Request, id string) error {
	task, err := DB.Tasks[id].Read(".html")
	if err != nil {
		return err
	}

	p := &ViewPage{
		Id:    id,
		Title: DB.Tasks[id].Title,
		Task:  string(task),
	}

	err = RenderTemplate(w, TemplatePath+"task.html", p)
	if err != nil {
		return err
	}

	return nil
}

// Handle requests on /plon/user/.
// List all of specified users tasks by rendering the user.html template.
func UserHandler(w http.ResponseWriter, r *http.Request, usernames []string) error {
	if usernames[0] != "wszyscy" {
		usernames = append(usernames, "wszyscy")
	}

	p := &UserPage{
		Username: usernames[0],
	}
	for _, username := range usernames {
		for id, _ := range DB.Users[username].Tasks {
			if DB.Users[username].Tasks[id] {
				p.Alltasks = append(p.Alltasks,
					struct {
						Id    string
						Title string
					}{
						Id:    id,
						Title: DB.Tasks[id].Title,
					})
			}
		}
	}

	err := RenderTemplate(w, TemplatePath+"user.html", p)
	if err != nil {
		return err
	}

	return nil
}
