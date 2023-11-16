package handlers

import (
	"io"
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
			c.search(w, email)
		}
	case http.MethodPost:
		c.addOne(w, r.Body)
	case http.MethodDelete:
		c.deleteOne(w, email)
	case http.MethodPut:
		c.updateOne(w, r.Body, email)
	default:
		http.Error(w, "ki korchish bhai... keno korchish erom", http.StatusMethodNotAllowed)
	}
}

func (c Contact) updateOne(w http.ResponseWriter, body io.Reader, email string) {
	contact := contact.Contact{}

	contact.FromJSON(body)

	cont, err := db.ToDoc(db.FindOne(c.l, email))

	if err != nil {
		c.l.Println(err)
		http.Error(w, "Error in fetching document from db", http.StatusInternalServerError)
	}

	_, err = c.db.Put(cont.Docs[0].Id, contact, cont.Docs[0].Rev)

	if err != nil {
		c.l.Println(err)
		http.Error(w, "Failed to update document in db", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
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

func (c Contact) addOne(w http.ResponseWriter, body io.Reader) {
	new := contact.Contact{}
	new.FromJSON(body)
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

	contacts.ToJSON(w)
}

func (c Contact) search(w http.ResponseWriter, email string) {
	dbResp := db.FindResp{}
	data := db.Search(c.l, email)

	dbResp.FromJSON(data)

	c.l.Println(dbResp.Docs)

	dbResp.Docs.ToJSON(w)
}

func NewContact(l *log.Logger, db *couchdb.DB) *Contact {
	return &Contact{
		l,
		db,
	}
}
