package db

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"text/template"
)

func Search(l *log.Logger, q string) io.Reader {

	query, _ := template.ParseFiles("./handlers/findQuery.txt")
	qbuf := bytes.Buffer{}
	query.Execute(&qbuf, map[string]string{
		"q": q,
	})

	l.Println(qbuf.String())

	resp, err := http.Post("http://admin:Tr0069er@localhost:5984/sdb_test/_find", "application/json", &qbuf)

	if err != nil {
		l.Println("Error in querying db", err)
		return nil
	}

	return resp.Body
}

func FindOne(l *log.Logger, email string) io.Reader {

	query, _ := template.ParseFiles("./handlers/searchQuery.txt")

	qbuf := bytes.Buffer{}

	query.Execute(&qbuf, map[string]string{
		"email": email,
	})

	resp, err := http.Post("http://admin:Tr0069er@localhost:5984/sdb_test/_find", "application/json", &qbuf)

	if err != nil {
		l.Println("Error in querying db", err)
		return nil
	}

	return resp.Body
}

func ToDoc(body io.Reader) (*DelResp, error) {

	if body == nil {
		return nil, errors.New("Nil bodyfer")
	}

	dbresp := DelResp{}

	err := dbresp.FromJSON(body)

	if err != nil {
		return nil, err
	}

	return &dbresp, nil

}

func ToContact(body io.Reader) (*FindResp, error) {

	if body == nil {
		return nil, errors.New("Nil bodyfer")
	}

	dbresp := FindResp{}

	err := dbresp.FromJSON(body)

	if err != nil {
		return nil, err
	}

	return &dbresp, nil
}
