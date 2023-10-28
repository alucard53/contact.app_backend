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

	// response := DBResp{}

	// err = db.AllDocs(&response, couchdb.Options{})

	// if err != nil {
	// 	l.Fatal("Error in query all docs", err)
	// }

	// l.Println(response.Rows[0].Id)

	// contact := contact.Contact{}

	// err = db.Get(response.Rows[0].Id, &contact, couchdb.Options{})

	// if err != nil {
	// 	l.Fatal("Failed to fetch object with key", err)
	// }

	// l.Println("Doc, ", contact)

	server := initServer(l, db)

	startServer(server)
}
