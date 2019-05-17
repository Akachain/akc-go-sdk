package akchtc

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SampleChaincode struct {
}

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error

	if fn == "insert" {
		result, err = insert(stub, args)
	} else if fn == "get" {
		result, err = get(stub, args)
	} else if fn == "prune" {
		result, err = prune(stub, args)
	} else if fn == "delete" {
		result, err = delete(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Demo use akachain high throught put in chaincode
func insert(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	akc := AkcHighThroughput{}
	check := akc.Insert(stub, args)

	return fmt.Sprintf("%s", check), check
}

func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	akc := AkcHighThroughput{}
	check, err := akc.Get(stub, args)

	return fmt.Sprintf("%f", check), err
}

func prune(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	akc := AkcHighThroughput{}
	check, err := akc.Prune(stub, args)

	return strconv.FormatBool(check), err
}

func delete(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	akc := AkcHighThroughput{}
	check, err := akc.Delete(stub, args)

	return strconv.FormatBool(check), err
}

func main() {
	err := shim.Start(new(SampleChaincode))
	if err != nil {
		fmt.Println("Could not start SampleChaincode")
	} else {
		fmt.Println("SampleChaincode successfully started")
	}
}
