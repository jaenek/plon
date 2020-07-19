package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

var DB Database

// TODO(#3): Add database creation and backup.
func main() {
	err := DB.Read("db.json")
	if err != nil {
		// Create empty database and send a warning instead of fatal error.
		log.Fatal(err.Error())
	}

	http.HandleFunc("/plon/add/", AddHandler)
	http.HandleFunc("/plon/edit/", MakeIdHandler(EditHandler))
	http.HandleFunc("/plon/save/", MakeIdHandler(SaveHandler))
	http.HandleFunc("/plon/delete/", MakeIdHandler(DeleteHandler))
	http.HandleFunc("/plon/view/", MakeIdHandler(ViewHandler))
	http.HandleFunc("/plon/user/", MakeUserHandler(UserHandler))
	http.HandleFunc("/plon/", MakeIndexHandler())
	http.HandleFunc("/plon", MakeIndexHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
