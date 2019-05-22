package util

import (
	"container/list"
	"encoding/json"
	"regexp"

	. "github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/op/go-logging"
)

// Logger for the shim package.
var mockLogger = logging.MustGetLogger("mockStubExtend")

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

	// A list containing ledger keys filtered by the query string
	var filteredByKeys = list.New()

	// Now go through all of the selector and try to filter
	// the State Map in order of the selector
	for key, value := range selector {
		// Filter state by key first
		if key == "_id" {
			switch vv := value.(type) {
			case string:
				// This is the case where we need to find the exact state key
				// This should never happen as the query should use stub.GetState directly instead of using this function
				// Anyway, we can just push the query value to the key filter list
				filteredByKeys.PushBack(vv)
			case interface{}:
				// This is the case where we need to find keys that satisfy specific regex
				// Our value is {"$regex": "^Quorum_"}
				i := vv.(map[string]interface{})
				regex := i["$regex"].(string)
				for j := stub.Keys.Front(); j != nil; j = j.Next() {
					// Loop through all the keys and find keys that match our regex
					kk := j.Value.(string)
					matched, _ := regexp.Match(regex, []byte(kk))
					if matched {
						// Found the matching key, push it into the list
						filteredByKeys.PushBack(kk)
					}
				}
			default:
				mockLogger.Error("Not implemented")
			}
		}
	}

	var filteredByAttributes = list.New()

	// Now go through all of the selector and try to filter
	// the filtered key list in order of the selector
	// Unfortunately, it is a n^2 loop
	for k, v := range selector {
		if k == "_id" {
			// we already filter first round above
			continue
		}

		for i := filteredByKeys.Front(); i != nil; i = i.Next() {
			// Get the state value
			stateValue, _ := stub.GetState(i.Value.(string))
			var item map[string]interface{}
			json.Unmarshal([]byte(stateValue), &item)

			// Loop through the JSON state to find matching attribute
			for jsonKey, jsonValue := range item {
				// Found matching attribute with the query selection
				if k == jsonKey {
					// Check
					switch vv := v.(type) {
					case string:
						if jsonValue.(string) == vv {
							filteredByAttributes.PushBack(i)
						}
					case interface{}:
						i := vv.(map[string]interface{})
						regex := i["$regex"].(string)
						matched, _ := regexp.Match(regex, []byte(jsonValue.(string)))
						if matched {
							filteredByAttributes.PushBack(i)
						}
					default:
						mockLogger.Error("Not implemented")
					}
				}
			}
		}

	}
	r := NewMockFilterQueryIterator(stub, filteredByAttributes)
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
