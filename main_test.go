package main

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/assert"
)

func setupMock() *util.MockStubExtend {
	// Initialize MockStubExtend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("sample", cc), cc)

	// Create a new database, Drop old database
	db, _ := util.NewCouchDBHandlerWithConnectionAuthentication(true)
	stub.SetCouchDBConfiguration(db)
	return stub
}

func TestUpdateData(t *testing.T) {
	stub := setupMock()
	key1 := "key1"
	key2 := "key2"
	val1 := "val1"
	val2 := "val2"

	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateData"), []byte(key1), []byte(key2), []byte("val0"), []byte("val0")})

	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("UpdateData"), []byte(key1), []byte(key2), []byte(val1), []byte(val2)})

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
}

func TestPartialQuery(t *testing.T) {
	stub := setupMock()
	key1 := "key1"
	val1 := "val1"
	val2 := "val2"

	// create 0=9 transactions sharing a part of the key "key_{number}"
	for i := 0; i < 9; i++ {
		util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateData"), []byte(key1), []byte(strconv.Itoa(i)), []byte(val1), []byte(val2)})
	}

	// Test GetStateByPartialCompositeKeyWithPagination
	resultsIterator, _ := stub.GetStateByPartialCompositeKey(DATATABLE, []string{key1})

	i := 0
	for resultsIterator.HasNext() {
		resultsIterator.Next()
		i++
	}

	assert.Equal(t, i, 9)
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

	// Prepare query string
	var queryString = `
	{ "selector": 
		{ 	
			"_id": 
				{"$gt": "\u0000Data_"}			
		}
	}`
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

func TestGetQueryResultWithPagination(t *testing.T) {
	stub := setupMock()
	keyPrefix := "key"
	valPrefix := "val"

	// Create 0-9 states with format "key_{number}" "val_{number}"
	for i := 0; i < 9; i++ {
		util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateData"), []byte(keyPrefix), []byte(strconv.Itoa(i)), []byte(valPrefix), []byte(strconv.Itoa(i))})
	}

	// Prepare query string
	var queryString = `
	{ "selector": 
		{ 	
			"_id": 
				{"$gt": "\u0000Data_\u0000key"}			
		}
	}`

	// fetch the first page with only 5
	var pageSize int32
	pageSize = 5
	_, queryResponse, _ := stub.GetQueryResultWithPagination(queryString, pageSize, "")
	assert.Equal(t, queryResponse.GetFetchedRecordsCount(), int32(5))

	// Now get the rest
	data, resp, _ := stub.GetQueryResultWithPagination(queryString, pageSize, queryResponse.Bookmark)
	assert.Equal(t, resp.GetFetchedRecordsCount(), int32(4))

	v, _ := data.Next()
	var dat Data
	json.Unmarshal(v.Value, &dat)
	assert.Equal(t, dat.Key2, "5")
}
