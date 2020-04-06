package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

var DB Database

func main() {
	err := DB.Read("db.json")
	if err != nil {
		// TODO: Create empty database and send a warning instead of fatal error.
		log.Fatal(err.Error())
	}

	http.HandleFunc("/plon/add/", AddHandler)
	http.HandleFunc("/plon/save/", SaveHandler)
	http.HandleFunc("/plon/view/", ViewHandler)
	http.HandleFunc("/plon/", IndexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
