package main

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

type User struct {
	Tasks map[string]bool
}

type Database struct {
	Users map[string]User
	Tasks map[string]Task
}

// Read database from json file.
func (db *Database) Read(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"file": path,
	}).Info("Unmarshalling json.")

	err = json.Unmarshal(b, db)
	if err != nil {
		return err
	}

	return nil
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
