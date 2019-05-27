package util

import (
	"container/list"
	"errors"
	"strings"

	. "github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	logging "github.com/op/go-logging"
)

// Logger for the shim package.
var mockLogger = logging.MustGetLogger("mockStubExtend")

// MockStubExtend provides composition class for MockStub as some of the mockstub methods are not implemented
type MockStubExtend struct {
	args      [][]byte        // this is private in MockStub
	cc        Chaincode       // this is private in MockStub
	CouchDB   bool            // if we use couchDB
	DbHandler *CouchDBHandler // if we use couchDB
	*MockStub
}

// GetQueryResult overrides the same function in MockStub
// that did not implement anything.
func (stub *MockStubExtend) GetQueryResult(query string) (StateQueryIteratorInterface, error) {
	// Query data from couchDB
	rawdata, _ := stub.DbHandler.QueryDocument(query)
	mockLogger.Debug(rawdata)

	// A list containing ledger keys filtered by the query string
	var filteredKeys = list.New()
	for _, k := range rawdata {
		//[]map[string]interface{}
		filteredKeys.PushBack(k["_id"])
	}

	// Test
	r := NewMockFilterQueryIterator(stub, filteredKeys)
	return r, nil
}

type MockStateQueryIterator struct {
	Closed  bool
	Data    *map[string][]byte
	Current *list.Element
}

func NewMockStubExtend(stub *MockStub, c Chaincode) *MockStubExtend {
	s := new(MockStubExtend)
	s.MockStub = stub
	s.cc = c
	s.CouchDB = false
	return s
}

func (stub *MockStubExtend) SetCouchDBConfiguration(handler *CouchDBHandler) {
	stub.CouchDB = true
	stub.DbHandler = handler
}

// Override this function from MockStub
func (stub *MockStubExtend) MockInvoke(uuid string, args [][]byte) pb.Response {
	stub.args = args
	stub.MockTransactionStart(uuid)
	res := stub.cc.Invoke(stub)
	stub.MockTransactionEnd(uuid)
	return res
}

// Override this function from MockStub
func (stub *MockStubExtend) MockInit(uuid string, args [][]byte) pb.Response {
	stub.args = args
	stub.MockTransactionStart(uuid)
	res := stub.cc.Init(stub)
	stub.MockTransactionEnd(uuid)
	return res
}

// Override this function from MockStub
func (stub *MockStubExtend) GetFunctionAndParameters() (function string, params []string) {
	allargs := stub.GetStringArgs()
	function = ""
	params = []string{}
	if len(allargs) >= 1 {
		function = allargs[0]
		params = allargs[1:]
	}
	return
}

// Override this function from MockStub
func (stub *MockStubExtend) GetStringArgs() []string {
	strargs := make([]string, 0, len(stub.args))
	for _, barg := range stub.args {
		strargs = append(strargs, string(barg))
	}
	return strargs
}

// PutState writes the specified `value` and `key` into the ledger.
func (stub *MockStubExtend) PutState(key string, value []byte) error {
	// Carry on
	stub.putStateOriginal(key, value)

	// In case we are using CouchDB, we will also store the value document in the database
	if stub.CouchDB {
		stub.DbHandler.SaveDocument(key, value)
	}

	return nil
}

// This is copied from mockstub as we still need to carry on normal putstate operation with the mock ledger map
func (stub *MockStubExtend) putStateOriginal(key string, value []byte) error {
	if stub.TxID == "" {
		err := errors.New("cannot PutState without a transactions - call stub.MockTransactionStart()?")
		mockLogger.Errorf("%+v", err)
		return err
	}

	// If the value is nil or empty, delete the key
	if len(value) == 0 {
		mockLogger.Debug("MockStub", stub.Name, "PutState called, but value is nil or empty. Delete ", key)
		return stub.DelState(key)
	}

	mockLogger.Debug("MockStub", stub.Name, "Putting", key, value)
	stub.State[key] = value

	// insert key into ordered list of keys
	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		elemValue := elem.Value.(string)
		comp := strings.Compare(key, elemValue)
		mockLogger.Debug("MockStub", stub.Name, "Compared", key, elemValue, " and got ", comp)
		if comp < 0 {
			// key < elem, insert it before elem
			stub.Keys.InsertBefore(key, elem)
			mockLogger.Debug("MockStub", stub.Name, "Key", key, " inserted before", elem.Value)
			break
		} else if comp == 0 {
			// keys exists, no need to change
			mockLogger.Debug("MockStub", stub.Name, "Key", key, "already in State")
			break
		} else { // comp > 0
			// key > elem, keep looking unless this is the end of the list
			if elem.Next() == nil {
				stub.Keys.PushBack(key)
				mockLogger.Debug("MockStub", stub.Name, "Key", key, "appended")
				break
			}
		}
	}

	// special case for empty Keys list
	if stub.Keys.Len() == 0 {
		stub.Keys.PushFront(key)
		mockLogger.Debug("MockStub", stub.Name, "Key", key, "is first element in list")
	}

	return nil
}
