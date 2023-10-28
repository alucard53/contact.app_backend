package db

import (
	"errors"
	"log"
	"net/http"

	"contact.app_backend/contact"
	"github.com/fjl/go-couchdb"
)

type Doc struct {
	Id string `json:"id"`
}

type DBResp struct {
	Rows []Doc `json:"rows"`
}

type FindResp struct {
	Docs contact.Contacts `json:"docs"`
}

func ConnectDB(url string, l *log.Logger) (*couchdb.DB, error) {
	client, err := couchdb.NewClient(url, http.DefaultTransport)

	if err != nil {
		l.Println("could not connect to db")
		return nil, err
	}

	l.Println("Connected to db server")

	db := client.DB("sdb_test")

	if db == nil {
		return nil, errors.New("DB not found")
	}

	return db, nil
}
