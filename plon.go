package main

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var DB Database

func main() {
	var backupPath string
	var export bool

	const (
		importDefault = ""
		importUsage = "Path to a backup file for example backup-202007202122.json."
		exportDefault = false
		exportUsage = "Export the database to a file with current date."
	)

	flag.StringVar(&backupPath, "import-db", importDefault, importUsage)
	flag.StringVar(&backupPath, "i", importDefault, importUsage+" (shorthand)")
	flag.BoolVar(&export, "export-db", exportDefault, exportUsage)
	flag.BoolVar(&export, "e", exportDefault, exportUsage+" (shorthand)")

	flag.Parse()

	if backupPath != "" {
		err := DB.Import(backupPath)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	err := ReadJSON("db.json", &DB)
	if err != nil {
		log.Warn(err.Error())

		DB, err = CreateDatabase("db.json")
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if export {
		err := DB.Export()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	http.HandleFunc("/plon/add/", AddHandler)
	http.HandleFunc("/plon/edit/", MakeIdHandler(EditHandler))
	http.HandleFunc("/plon/save/", MakeIdHandler(SaveHandler))
	http.HandleFunc("/plon/delete/", MakeIdHandler(DeleteHandler))
	http.HandleFunc("/plon/view/", MakeIdHandler(ViewHandler))
	http.HandleFunc("/plon/user/", MakeUserHandler(UserHandler))
	http.HandleFunc("/plon/", MakeIndexHandler())
	http.HandleFunc("/plon", MakeIndexHandler())
	log.Fatal(http.ListenAndServe(":8001", nil))
}
