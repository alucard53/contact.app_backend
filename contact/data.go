package contact

import (
	"encoding/json"
	"io"
)

type Contact struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     int    `json:"phone"`
	Email     string `json:"email"`
}

type Contacts []Contact

func (c *Contact) FromJSON(r io.Reader) {
	json.NewDecoder(r).Decode(c)
}

func (c Contacts) TOJSON(w io.Writer) {
	json.NewEncoder(w).Encode(c)
}
