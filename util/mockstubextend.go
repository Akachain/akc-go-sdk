package util

import (
	"errors"
	"fmt"

	. "github.com/hyperledger/fabric/core/chaincode/shim"
)

type MockStubExtendInterface interface {
	GetQueryResult(query string) StateQueryIteratorInterface
}

// MockStubExtend provides
type MockStubExtend struct {
	*MockStub
}

func (stub MockStubExtend) GetQueryResult(query string) (StateQueryIteratorInterface, error) {
	fmt.Println("MockStubExtend")
	return nil, errors.New("not implemented")
}

func NewMockStubExtend(stub *MockStub) *MockStubExtend {
	s := new(MockStubExtend)
	s.MockStub = stub
	return s
}
