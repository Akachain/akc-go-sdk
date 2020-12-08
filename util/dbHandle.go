package util

import (
	"encoding/json"
	"github.com/hyperledger/fabric/common/metrics/disabled"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/statecouchdb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/version"
	"github.com/hyperledger/fabric/core/ledger/util/couchdb"
	//"path/filepath"

	//"github.com/hyperledger/fabric/core/ledger"
	"io/ioutil"
	"time"
)

const (
	// DefaultBaseURL is the default address of CouchDB server.
	DefaultBaseURL = "localhost:5984"

	// The couchDB test will have this name: DefaultChannelName_DefaultNamespace
	DefaultChannelName   = "channel"   // Fabric channel
	DefaultChaincodeName = "chaincode" // Fabric chaincode
)

// CouchDBHandler holds 1 parameter:
// dbEngine: a VersionedDB object that is used by the chaincode to query.
// This is to guarantee that the test uses the same logic in interaction with stateDB as the chaincode.
// This also includes how chaincode builds its query to interact with the stateDB.
type CouchDBHandler struct {
	dbEngine *statecouchdb.VersionedDB
}

// NewCouchDBHandlerWithConnectionAuthentication returns a new CouchDBHandler and setup database for testing
func NewCouchDBHandlerWithConnectionAuthentication(isDrop bool) (*CouchDBHandler, error) {
	// Sometimes we'll have to drop the database to clean all previous test
	if isDrop == true {
		cleanUp()
	}

	// Create a new dbEngine for the channel
	handler := new(CouchDBHandler)
	redoPath, err := ioutil.TempDir("", "redoPath")
	config := &couchdb.Config{
		Address:             DefaultBaseURL,
		InternalQueryLimit:  1000,
		MaxBatchUpdateSize:  1000,
		MaxRetries:          3,
		MaxRetriesOnStartup: 20,
		RequestTimeout:      35 * time.Second,
		RedoLogPath:         redoPath,
		UserCacheSizeMBs:    8,
	}
	couchState, _ := statecouchdb.NewVersionedDBProvider(config, &disabled.Provider{}, &statedb.Cache{})

	// This step creates a redundant meta database with name channel_ ,
	// there should be some ways to prevent this. We leave it for now
	h, err := couchState.GetDBHandle(DefaultChannelName)
	if err != nil {
		return nil, err
	}
	handler.dbEngine = h.(*statecouchdb.VersionedDB)
	return handler, nil
}

func cleanUp() error {
	// statedb.VersionedDB does not publish its couchDB object
	// Thus, we'll have to recreate
	// set required config data to use state couchdb
	redoPath, _ := ioutil.TempDir("", "redoPath")
	couchdbConfig := &couchdb.Config{
		Address:             DefaultBaseURL,
		Username:            "",
		Password:            "",
		MaxRetries:          3,
		MaxRetriesOnStartup: 20,
		RequestTimeout:      35 * time.Second,
		RedoLogPath:         redoPath,
	}
	ins, er := couchdb.CreateCouchInstance(couchdbConfig, &disabled.Provider{})
	if er != nil {
		return er
	}
	dbName := couchdb.ConstructNamespaceDBName(DefaultChannelName, DefaultChaincodeName)
	db := couchdb.CouchDatabase{CouchInstance: ins, DBName: dbName}
	_, er = db.DropDatabase()
	return er
}

// NewCouchDBHandlerWithConnection that is compatibles with previous release
func NewCouchDBHandlerWithConnection(dbName string, isDrop bool, connectionString string) (*CouchDBHandler, error) {
	return NewCouchDBHandlerWithConnectionAuthentication(isDrop)
}

// NewCouchDBHandler that is compatibles with previous release
func NewCouchDBHandler(dbName string, isDrop bool) (*CouchDBHandler, error) {
	return NewCouchDBHandlerWithConnection(dbName, isDrop, DefaultBaseURL)
}

// SaveDocument stores a value in couchDB
func (handler *CouchDBHandler) SaveDocument(key string, value []byte) error {
	// unmarshal the value param
	var doc map[string]interface{}
	json.Unmarshal(value, &doc)

	// Save the doc in database
	batch := statedb.NewUpdateBatch()
	batch.Put(DefaultChaincodeName, key, value, version.NewHeight(1, 1))
	savePoint := version.NewHeight(1, 2)
	err := handler.dbEngine.ApplyUpdates(batch, savePoint)

	return err
}

// QueryDocument executes a query string and return results
func (handler *CouchDBHandler) QueryDocument(query string) (statedb.ResultsIterator, error) {
	rs, er := handler.dbEngine.ExecuteQuery(DefaultChaincodeName, query)
	return rs, er
}

// QueryDocumentWithPagination executes a query string and return results
func (handler *CouchDBHandler) QueryDocumentWithPagination(query string, limit int32, bookmark string) (statedb.ResultsIterator, error) {
	queryOptions := make(map[string]interface{})
	if limit != 0 {
		queryOptions["limit"] = limit
	}
	if bookmark != "" {
		queryOptions["bookmark"] = bookmark
	}
	rs, er := handler.dbEngine.ExecuteQueryWithMetadata(DefaultChaincodeName, query, queryOptions)
	return rs, er
}

// ReadDocument executes a query string and return results
func (handler *CouchDBHandler) ReadDocument(id string) ([]byte, error) {
	rs, er := handler.dbEngine.GetState(DefaultChaincodeName, id)
	if er != nil {
		return nil, er
	}
	// found no document in db with id
	if rs == nil {
		return nil, nil
	}
	return rs.Value, er
}

// QueryDocumentByRange get a list of documents from couchDB by key range
func (handler *CouchDBHandler) QueryDocumentByRange(startKey, endKey string) (statedb.ResultsIterator, error) {
	rs, er := handler.dbEngine.GetStateRangeScanIterator(DefaultChaincodeName, startKey, endKey)
	return rs, er
}

//// QueryDocumentByRange get a list of documents from couchDB by key range
//// TODO: GetStateRangeScanIteratorWithMetadata does not accept bookmark
//func (handler *CouchDBHandler) QueryDocumentByRangeWithPagination(startKey, endKey string, limit int32, bookmark string) (statedb.ResultsIterator, error) {
//	queryOptions := make(map[string]interface{})
//	if limit != 0 {
//		queryOptions["limit"] = limit
//	}
//	//if bookmark != "" {
//	//	queryOptions["bookmark"] = bookmark
//	//}
//
//	rs, er := handler.dbEngine.GetStateRangeScanIteratorWithMetadata(DefaultChaincodeName, startKey, endKey, queryOptions)
//	return rs, er
//}
