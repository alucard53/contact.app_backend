package db

import "contact.app_backend/contact"

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
