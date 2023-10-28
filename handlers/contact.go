package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"text/template"

	"contact.app_backend/contact"
	"contact.app_backend/db"
	"github.com/fjl/go-couchdb"
	"github.com/google/uuid"
)

type Contact struct {
	l  *log.Logger
	db *couchdb.DB
}

func (c Contact) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		email := r.URL.Query().Get("e")
		if email == "" {
			c.getAll(w)
		} else {
			c.findOne(w, email)
		}
	case http.MethodPost:
		c.addOne(w, r)
	default:
	}
}

func (c Contact) addOne(w http.ResponseWriter, r *http.Request) {
	new := contact.Contact{}
	new.FromJSON(r.Body)
	c.l.Println(new)
	_, err := c.db.Put(uuid.NewString(), new, "")

	if err != nil {
		c.l.Println(err)
		http.Error(w, "Error in inserting document", http.StatusInternalServerError)
		return
	}
}

func (c Contact) getAll(w http.ResponseWriter) {
	response := db.DBResp{}

	err := c.db.AllDocs(&response, couchdb.Options{})

	if err != nil {
		c.l.Println("Error in querying docs", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contacts := make(contact.Contacts, len(response.Rows))

	for i, row := range response.Rows {
		err = c.db.Get(row.Id, &contacts[i], couchdb.Options{})
		if err != nil {
			c.l.Println("Error in querying docs", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	contacts.TOJSON(w)

}

func (c Contact) findOne(w http.ResponseWriter, email string) {
	query, err := template.ParseFiles("./handlers/searchQuery.txt")

	if err != nil {
		http.Error(w, "Error in reading query file", http.StatusInternalServerError)
		return
	}

	qbuf := bytes.Buffer{}

	query.Execute(&qbuf, map[string]string{
		"email": email,
	})

	resp, err := http.Post("http://admin:Tr0069er@localhost:5984/sdb_test/_find", "application/json", &qbuf)

	if err != nil {
		http.Error(w, "Error in fetchingn query from db", http.StatusInternalServerError)
		return
	}

	rbuf, _ := io.ReadAll(resp.Body)

	dbresp := db.FindResp{}

	err = json.Unmarshal(rbuf, &dbresp)

	if err != nil {
		c.l.Println(err)
		return
	}

	dbresp.Docs.TOJSON(w)
}

func NewContact(l *log.Logger, db *couchdb.DB) *Contact {
	return &Contact{
		l,
		db,
	}
}
