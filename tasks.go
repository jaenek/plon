package main

import (
	"io/ioutil"
	"os"
	"path"
	"time"
)

// Declare directory where tasks are stored.
const TaskPath = "tasks/"

type Task struct {
	Title      string
	Path       string
	Addressees []string
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

	err := ioutil.WriteFile(t.Path, []byte(task), 0644)
	if err != nil {
		return err
	}

	for _, username := range t.Addressees {
		user := DB.Users[username]
		user.Tasks = append(DB.Users[username].Tasks, id)
		DB.Users[username] = user
	}

	DB.Tasks[id] = t
	DB.Write()

	return nil
}

// Read the task from file.
func (t Task) Read() (string, error) {
	task, err := ioutil.ReadFile(t.Path)
	if err != nil {
		return "", err
	}
	return string(task), nil
}
