package util

import (
	"errors"

	. "github.com/hyperledger/fabric/core/chaincode/shim"
	couchdb "github.com/hyperledger/fabric/core/ledger/util/couchdb"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
)

// AkcQueryIterator inherits StateQueryIterator to simulate how the peer handle query string response
type AkcQueryIterator struct {
	data       []*couchdb.QueryResult
	currentLoc int
	*StateQueryIterator
}

func (it *AkcQueryIterator) HasNext() bool {
	return it.currentLoc < len(it.data)
}

func (it *AkcQueryIterator) Next() (*queryresult.KV, error) {
	var kv = new(queryresult.KV)

	if !it.HasNext() {
		return nil, errors.New("There is no other item in the iterator")
	}

	item := it.data[it.currentLoc]
	it.currentLoc++

	if item == nil {
		return nil, errors.New("Empty query result")
	}

	kv.Value = item.Value

	return kv, nil
}

func (iter *AkcQueryIterator) Close() error {
	return nil
}
