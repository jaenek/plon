package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
)

type User struct {
	Tasks map[string]bool
}

type Database struct {
	Users map[string]User
	Tasks map[string]Task
}

type TaskBackup struct {
	Title      string
	Content    string
	Addressees map[string]bool
	Created    time.Time
	Due        time.Time
}

type Backup struct {
	Users map[string]User
	Tasks map[string]TaskBackup
}

func (db *Database) Import(path string) error {
	b := &Backup{}

	err := ReadJSON(path, b)
	if err != nil {
		return err
	}

	db = &Database{
		Users: b.Users,
		Tasks: make(map[string]Task, len(b.Tasks)),
	}

	for id, task := range b.Tasks {
		db.Tasks[id] = Task{
			Title: task.Title,
			Path: TaskPath + id + "/" + id,
			Addressees: task.Addressees,
			Created: task.Created,
			Due: task.Due,
		}

		err = db.Tasks[id].Save(id, task.Content)
		if err != nil {
			return err
		}
	}

	err = WriteJSON("db.json", db)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Export() error {
	b := &Backup{
		Users: db.Users,
		Tasks: make(map[string]TaskBackup, len(db.Tasks)),
	}

	for id, task := range db.Tasks {
		content, err := ioutil.ReadFile(task.Path+".md")
		if err != nil {
			return err
		}

		b.Tasks[id] = TaskBackup{
			Title: task.Title,
			Content: string(content),
			Addressees: task.Addressees,
			Created: task.Created,
			Due: task.Due,
		}
	}

	t := time.Now()
	err := WriteJSON("backup-"+t.Format("200601021504")+".json", b)
	if err != nil {
		return err
	}

	return nil
}

// Read database from json file.
func ReadJSON(path string, data interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"file": path,
	}).Info("Unmarshalling json.")

	err = json.Unmarshal(b, data)
	if err != nil {
		return err
	}

	return nil
}

// Write database to json file.
func WriteJSON(path string, data interface{}) error {
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"file": path,
	}).Info("Saving to file.")
	return nil
}
