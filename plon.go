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
	http.HandleFunc("/plon/edit/", MakeHandler(EditHandler))
	http.HandleFunc("/plon/save/", MakeHandler(SaveHandler))
	http.HandleFunc("/plon/delete/", MakeHandler(DeleteHandler))
	http.HandleFunc("/plon/view/", MakeHandler(ViewHandler))
	http.HandleFunc("/plon/", MakeIndexHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
