package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/fjl/go-couchdb"
)

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

func FindOne(l *log.Logger, email string) []byte {

	query, _ := template.ParseFiles("./handlers/searchQuery.txt")

	qbuf := bytes.Buffer{}

	query.Execute(&qbuf, map[string]string{
		"email": email,
	})

	resp, err := http.Post("http://admin:Tr0069er@localhost:5984/sdb_test/_find", "application/json", &qbuf)

	if err != nil {
		l.Println("Error in querying db")
		return nil
	}

	rbuf, err := io.ReadAll(resp.Body)

	l.Printf("%s", rbuf)

	if err != nil {
		l.Println("Error reading response body")
		return nil
	}

	return rbuf
}

func ToDoc(buf []byte) (*DelResp, error) {

	if buf == nil {
		return nil, errors.New("Nil buffer")
	}

	dbresp := DelResp{}

	err := json.Unmarshal(buf, &dbresp)

	if err != nil {
		return nil, err
	}
	return &dbresp, nil

}

func ToContact(buf []byte) (*FindResp, error) {

	if buf == nil {
		return nil, errors.New("Nil buffer")
	}

	dbresp := FindResp{}

	err := json.Unmarshal(buf, &dbresp)

	if err != nil {
		return nil, err
	}
	return &dbresp, nil
}
