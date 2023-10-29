package handlers

import (
	"log"
	"net/http"

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

	email := r.URL.Query().Get("e")
	switch r.Method {
	case http.MethodGet:
		if email == "" {
			c.getAll(w)
		} else {
			c.findOne(w, email)
		}
	case http.MethodPost:
		c.addOne(w, r)
	case http.MethodDelete:
		c.deleteOne(w, email)
	default:
		http.Error(w, "ki korchish bhai... keno korchish erom", http.StatusMethodNotAllowed)
	}
}

func (c Contact) deleteOne(w http.ResponseWriter, email string) {

	cont, err := db.ToDoc(db.FindOne(c.l, email))

	if err != nil || len(cont.Docs) < 1 {
		c.l.Println(err)
		http.Error(w, "Error in deleting document", http.StatusInternalServerError)
		return
	}

	c.l.Println(cont)

	_, err = c.db.Delete(cont.Docs[0].Id, cont.Docs[0].Rev)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	cont, err := db.ToContact(db.FindOne(c.l, email))

	if err != nil {
		c.l.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cont.Docs.TOJSON(w)
}

func NewContact(l *log.Logger, db *couchdb.DB) *Contact {
	return &Contact{
		l,
		db,
	}
}
