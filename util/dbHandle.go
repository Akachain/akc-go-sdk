package util

import (
	"encoding/json"
	"time"

	"github.com/hyperledger/fabric/common/metrics/disabled"
	couchdb "github.com/hyperledger/fabric/core/ledger/util/couchdb"
)

const (
	// DefaultBaseURL is the default address of CouchDB server.
	DefaultBaseURL = "localhost:5984"
)

type CouchDBHandler struct {
	CouchDatabase *couchdb.CouchDatabase
}

// NewCouchDBHandlerWithConnection returns a new CouchDBHandler and setup database for testing
func NewCouchDBHandlerWithConnection(dbName string, isDrop bool, connectionString string) (*CouchDBHandler, error) {
	handler := new(CouchDBHandler)

	//Create a couchdb instance
	couchDBInstance, er := couchdb.CreateCouchInstance(connectionString, "", "", 3, 10, time.Second*30, true, &disabled.Provider{})
	if er != nil {
		return nil, er
	}

	//Create a couchdatabase
	db := couchdb.CouchDatabase{CouchInstance: couchDBInstance, DBName: dbName}
	if isDrop == true {
		db.DropDatabase()
	}

	er = db.CreateDatabaseIfNotExist()
	if er != nil {
		return nil, er
	}

	handler.CouchDatabase = &db
	return handler, nil
}

// NewCouchDBHandler returns a new CouchDBHandler and setup database for testing
func NewCouchDBHandler(dbName string, isDrop bool) (*CouchDBHandler, error) {
	return NewCouchDBHandlerWithConnection(dbName, isDrop, DefaultBaseURL)
}

// SaveDocument stores a value in couchDB
func (handler *CouchDBHandler) SaveDocument(key string, value []byte) (string, error) {
	// unmarshal the value param
	var doc map[string]interface{}
	json.Unmarshal(value, &doc)

	// Save the doc in database
	rev, err := handler.CouchDatabase.SaveDoc(key, "", &couchdb.CouchDoc{JSONValue: value, Attachments: nil})
	return rev, err
}

// UpdateDocument update a value in couchDB
func (handler *CouchDBHandler) UpdateDocument(key string, value []byte) error {
	// unmarshal the value param
	var doc map[string]interface{}
	json.Unmarshal(value, &doc)

	_, rev, _ := handler.CouchDatabase.ReadDoc(key)

	// Save the doc in database
	_, err := handler.CouchDatabase.SaveDoc(key, rev, &couchdb.CouchDoc{JSONValue: value, Attachments: nil})
	return err
}

// QueryDocument executes a query string and return results
func (handler *CouchDBHandler) QueryDocument(query string) ([]*couchdb.QueryResult, error) {
	rs, _, er := handler.CouchDatabase.QueryDocuments(query)
	return rs, er
}

// ReadDocument executes a query string and return results
func (handler *CouchDBHandler) ReadDocument(id string) (*couchdb.CouchDoc, string, error) {
	couchDoc, revision, er := handler.CouchDatabase.ReadDoc(id)
	return couchDoc, revision, er
}

// DeleteDocument delete a specific document on couchdb
func (handler *CouchDBHandler) DeleteDocument(key string) error {
	_, rev, _ := handler.CouchDatabase.ReadDoc(key)
	er := handler.CouchDatabase.DeleteDoc(key, rev)
	return er
}

// QueryDocumentByRange get a list of documents from couchDB by key range
func (handler *CouchDBHandler) QueryDocumentByRange(startKey, endKey string, limit int32) ([]*couchdb.QueryResult, string, error) {
	rs, next, er := handler.CouchDatabase.ReadDocRange(startKey, endKey, limit)
	return rs, next, er
}
