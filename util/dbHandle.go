package util

import "github.com/leesper/couchdb-golang"

const (
	// DefaultBaseURL is the default address of CouchDB server.
	DefaultBaseURL = "http://localhost:5984"
)

type CouchDBHandler struct {
	Database *couchdb.Database
}

func NewCouchDBHandler() *CouchDBHandler {
	s := new(CouchDBHandler)
	db, _ := couchdb.NewDatabase(DefaultBaseURL)
	s.Database = db
	return s
}
