package util

import (
	"encoding/json"

	couchdb "github.com/leesper/couchdb-golang"
)

const (
	// DefaultBaseURL is the default address of CouchDB server.
	DefaultBaseURL = "http://localhost:5984"
)

type CouchDBHandler struct {
	Database *couchdb.Database
}

// NewCouchDBHandler returns a new CouchDBHandler and setup database for testing
func NewCouchDBHandler(dbName string) *CouchDBHandler {
	handler := new(CouchDBHandler)
	server, _ := handler.SetupServer(DefaultBaseURL)
	handler.SetupDB(dbName, server)
	return handler
}

// SetupServer creates a new couchDB server instance
func (*CouchDBHandler) SetupServer(url string) (*couchdb.Server, error) {
	server, err := couchdb.NewServer(url)
	return server, err
}

// SetupDB creates a new database instance with specific name
func (handler *CouchDBHandler) SetupDB(name string, server *couchdb.Server) (*couchdb.Database, error) {
	server.Delete(name)
	var err error
	db, err := server.Create(name)
	handler.Database = db
	return db, err
}

// SaveDocument stores a value in couchDB
func (handler *CouchDBHandler) SaveDocument(key string, value []byte) (string, error) {
	// unmarshal the value param
	var doc map[string]interface{}
	json.Unmarshal(value, &doc)
	// Add key as document id
	doc["_id"] = key
	// Save the doc in database
	id, _, err := handler.Database.Save(doc, nil)
	return id, err
}

// QueryDocument executes a query string and return results
func (handler *CouchDBHandler) QueryDocument(query string) ([]map[string]interface{}, error) {
	docsRaw, err := handler.Database.QueryJSON(query)
	return docsRaw, err
}
