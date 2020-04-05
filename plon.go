package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	http.HandleFunc("/plon/view/", ViewHandler)
	http.HandleFunc("/plon/", IndexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
