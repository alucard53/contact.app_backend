package main

import (
	"log"
	"os"

	"contact.app_backend/db"
)

func main() {
	l := log.New(os.Stdout, "->", log.Default().Flags())

	db, err := db.ConnectDB("http://admin:Tr0069er@localhost:5984", l)

	if err != nil {
		l.Fatal(err)
	}

	server := initServer(l, db)

	startServer(server)
}
