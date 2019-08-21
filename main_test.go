package main

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	DbURL  string `json:"TEST_COUCHDB_URL"`
	DbName string `json:"TEST_DATABASE_NAME"`
}

func getTestConfig(fileName string) (testConfig, error) {
	// fileName is the path to the json config file
	file, err := os.Open(fileName)
	var cfg testConfig
	if err != nil {
		return cfg, err
	}

	// decode to get config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func setupMock() *util.MockStubExtend {
	// Fetch test configuration
	cfg, _ := getTestConfig("config.json")

	// Initialize mockstubextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("sample", cc), cc)

	// Create a new database, Drop old database
	db, _ := util.NewCouchDBHandlerWithConnection(cfg.DbName, true, cfg.DbURL)
	stub.SetCouchDBConfiguration(db)
	return stub
}

func TestPartialQuery(t *testing.T) {
	stub := setupMock()
	key1 := "key"
	val1 := "val1"
	val2 := "val2"

	// create 0=9 transactions sharing a part of the key "key_{number}"
	for i := 0; i < 9; i++ {
		util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateData"), []byte(key1), []byte(strconv.Itoa(i)), []byte(val1), []byte(val2)})
	}

	// Test GetStateByPartialCompositeKeyWithPagination
	resultsIterator, queryResponse, _ := stub.GetStateByPartialCompositeKeyWithPagination(DATATABLE, []string{key1}, 5, "")

	i := 0
	for resultsIterator.HasNext() {
		resultsIterator.Next()
		i++
	}

	assert.Equal(t, i, 5)

	resultsIterator, _, _ = stub.GetStateByPartialCompositeKeyWithPagination(DATATABLE, []string{key1}, 5, queryResponse.GetBookmark())

	i = 0
	for resultsIterator.HasNext() {
		resultsIterator.Next()
		i++
	}

	assert.Equal(t, i, 4)
}

func TestSimpleData(t *testing.T) {
	stub := setupMock()
	key1 := "key1"
	key2 := "key2"
	key3 := "key1'"
	key4 := "key2'"
	val1 := "val1"
	val2 := "val2"

	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateData"), []byte(key1), []byte(key2), []byte(val1), []byte(val2)})

	// Check if the created data exist in the ledger
	compositeKey, _ := stub.CreateCompositeKey(DATATABLE, []string{key1, key2})
	state, _ := stub.GetState(compositeKey)
	var ad [10]Data

	json.Unmarshal([]byte(state), &ad[0])

	// Check if the created data information is correct
	assert.Equal(t, key1, ad[0].Key1)
	assert.Equal(t, key2, ad[0].Key2)
	assert.Equal(t, val1, ad[0].Attribute1)
	assert.Equal(t, val2, ad[0].Attribute2)

	// Test query string
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateData"), []byte(key3), []byte(key4), []byte(val1), []byte(val2)})

	queryString := "{\"selector\": {\"_id\": {\"$regex\": \"Data_\"}}}"
	resultsIterator, _ := stub.GetQueryResult(queryString)

	i := 0
	for resultsIterator.HasNext() {
		queryResponse, _ := resultsIterator.Next()
		json.Unmarshal(queryResponse.Value, &ad[i])
		i++
	}

	// Check if the created data information is correct
	assert.Equal(t, key1, ad[0].Key1)
	assert.Equal(t, key2, ad[0].Key2)
	assert.Equal(t, val1, ad[0].Attribute1)
	assert.Equal(t, val2, ad[0].Attribute2)
	assert.Equal(t, key3, ad[1].Key1)
	assert.Equal(t, key4, ad[1].Key2)
}
