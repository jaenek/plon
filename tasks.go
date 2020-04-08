package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/russross/blackfriday/v2"
)

// Declare directory where tasks are stored.
const TaskPath = "tasks/"

type Task struct {
	Title      string
	Path       string
	Addressees map[string]bool
	Created    time.Time
	Due        time.Time
}

// Create task directory.
// Save task with specified id to file.
// Add the task to the database and save it.
func (t Task) Save(id string, task string) error {
	if _, err := os.Stat(path.Dir(t.Path)); os.IsNotExist(err) {
		err := os.Mkdir(path.Dir(t.Path), 0755)
		if err != nil {
			return err
		}
	}

	err := ioutil.WriteFile(t.Path+".md", []byte(task), 0644)
	if err != nil {
		return err
	}

	data := bytes.Replace([]byte(task), []byte("\r"), nil, -1)
	html := blackfriday.Run(data)

	err = ioutil.WriteFile(t.Path+".html", html, 0644)
	if err != nil {
		return err
	}

	for username := range DB.Users {
		if t.Addressees[username] {
			DB.Users[username].Tasks[id] = true
		} else {
			DB.Users[username].Tasks[id] = false
		}
	}

	DB.Tasks[id] = t
	DB.Write()

	return nil
}

// Read the task from file with specified extension.
func (t Task) Read(extension string) (string, error) {
	task, err := ioutil.ReadFile(t.Path + extension)
	if err != nil {
		return "", err
	}
	return string(task), nil
}

// Delete the task file.
func (t Task) Delete(id string) error {
	err := os.RemoveAll(path.Dir(t.Path))
	if err != nil {
		return err
	}

	for _, user := range DB.Users {
		delete(user.Tasks, id)
	}

	delete(DB.Tasks, id)
	DB.Write()

	return nil
}
