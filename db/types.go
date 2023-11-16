package db

import (
	"encoding/json"
	"io"

	"contact.app_backend/contact"
)

type DBDoc struct {
	Id  string `json:"_id"`
	Rev string `json:"_rev"`
}

type DBResp struct {
	Rows []struct {
		Id string `json:"id"`
	} `json:"rows"`
}

type FindResp struct {
	Docs contact.Contacts `json:"docs"`
}

type DelResp struct {
	Docs []DBDoc `json:"docs"`
}

func (dbr *FindResp) FromJSON(r io.Reader) error {
	err := json.NewDecoder(r).Decode(dbr)

	if err != nil {
		return err
	}

	return nil
}

func (dbr *DelResp) FromJSON(r io.Reader) error {
	err := json.NewDecoder(r).Decode(dbr)
	if err != nil {
		return err
	}
	return nil
}
