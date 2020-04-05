package main

import (
	"encoding/json"
	"io/ioutil"
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

func ReadDatabase(path string) (*Database, error) {
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

func (db *Database) Write() error {
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
