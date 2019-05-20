package util

import (
	"encoding/json"
	"fmt"

	. "github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// MockStubExtend provides composition class for MockStub
type MockStubExtend struct {
	args [][]byte  // this is private in MockStub
	cc   Chaincode // this is private in MockStub
	*MockStub
}

type q struct {
	regex string
}

// GetQueryResult overrides the same function in MockStub
// that did not implement anything.
func (stub *MockStubExtend) GetQueryResult(query string) (StateQueryIteratorInterface, error) {

	// A sample query string is like this
	// {"selector": {"_id": {"$regex": "^Quorum_"},"ProposalID": "a"}}
	// we unmarshall the query string into a map of interface.
	var result map[string]interface{}
	json.Unmarshal([]byte(query), &result)
	selector := result["selector"].(map[string]interface{})

	// Now go through all of the selector and try to filter
	// the State Map in order of the selector
	for key, value := range selector {
		switch vv := value.(type) {
		case string:
			// This is the normal case
			// Our value is a normal string "a"
			fmt.Println(key, value.(string))
		case interface{}:
			// This is the nested object case
			// Our value is {"$regex": "^Quorum_"}
			i := vv.(map[string]interface{})
			for k, v := range i {
				if k == "regex" {
					fmt.Println(k, v)
				}
			}
		default:
			fmt.Println("Not implemented")
		}
	}

	r, _ := stub.GetStateByRange("1", "2")

	return r, nil
}

func NewMockStubExtend(stub *MockStub, c Chaincode) *MockStubExtend {
	s := new(MockStubExtend)
	s.MockStub = stub
	s.cc = c
	return s
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
