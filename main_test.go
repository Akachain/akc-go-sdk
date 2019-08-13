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

	for i := 0; i < 10; i++ {
		util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateData"), []byte(key1), []byte(strconv.Itoa(i)), []byte(val1), []byte(val2)})
	}

	stub.GetStateByPartialCompositeKey(DATATABLE, []string{key1})

}

func TestSimpleData(t *testing.T) {
	stub := setupMock()
	key1 := "key1"
	key2 := "key2"
	val1 := "val1"
	val2 := "val2"

	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateData"), []byte(key1), []byte(key2), []byte(val1), []byte(val2)})

	// Check if the created data exist in the ledger
	compositeKey, _ := stub.CreateCompositeKey(DATATABLE, []string{key1, key2})
	state, _ := stub.GetState(compositeKey)
	var ad Data
	json.Unmarshal([]byte(state), &ad)

	// Check if the created data information is correct
	assert.Equal(t, key1, ad.Key1)
	assert.Equal(t, key2, ad.Key2)
	assert.Equal(t, val1, ad.Attribute1)
	assert.Equal(t, val2, ad.Attribute2)
}
