package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
)

type User struct {
	Tasks []string
}

type Task struct {
	Title   string
	Path    string
	Created time.Time
	Due     time.Time
}

type Database struct {
	Users map[string]User
	Tasks map[string]Task
}

// Read database from json file.
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

// Write database to json file.
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
