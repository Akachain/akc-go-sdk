package util

import (
	"container/list"
	"errors"

	"github.com/hyperledger/fabric/protos/ledger/queryresult"
)

type mockFilterQueryIterator struct {
	Closed       bool
	Stub         *MockStubExtend
	FilteredKeys *list.List
	Current      *list.Element
}

// HasNext returns true if the query iterator contains additional keys and values.
func (iter *mockFilterQueryIterator) HasNext() bool {
	if iter.Closed {
		// previously called Close()
		mockLogger.Debug("HasNext() but already closed")
		return false
	}

	if iter.Current == nil {
		mockLogger.Error("HasNext() couldn't get Current")
		return false
	}

	current := iter.Current

	if current.Next() != nil {
		return true
	}

	return false
}

// Next returns the next key and value in the range query iterator.
func (iter *mockFilterQueryIterator) Next() (*queryresult.KV, error) {
	if iter.Closed == true {
		err := errors.New("mockFilterQueryIterator.Next() called after Close()")
		mockLogger.Errorf("%+v", err)
		return nil, err
	}

	if iter.HasNext() == false {
		err := errors.New("mockFilterQueryIterator.Next() called when it does not HaveNext()")
		mockLogger.Errorf("%+v", err)
		return nil, err
	}

	for iter.Current != nil {
		key := iter.Current.Value.(string)
		value, err := iter.Stub.GetState(key)
		iter.Current = iter.Current.Next()
		return &queryresult.KV{Key: key, Value: value}, err
	}

	err := errors.New("mockFilterQueryIterator.Next() went past end of range")
	mockLogger.Errorf("%+v", err)
	return nil, err
}

// Close closes the range query iterator. This should be called when done
// reading from the iterator to free up resources.
func (iter *mockFilterQueryIterator) Close() error {
	if iter.Closed == true {
		err := errors.New("mockFilterQueryIterator.Close() called after Close()")
		mockLogger.Errorf("%+v", err)
		return err
	}

	iter.Closed = true
	return nil
}

func NewMockFilterQueryIterator(stub *MockStubExtend, keys *list.List) *mockFilterQueryIterator {
	iter := new(mockFilterQueryIterator)
	iter.Closed = false
	iter.Stub = stub
	iter.FilteredKeys = keys
	iter.Current = stub.Keys.Front()

	return iter
}
