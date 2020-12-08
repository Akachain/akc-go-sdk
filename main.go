package main

import (
	"fmt"

	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// Chaincode implementation
type Chaincode struct {
}

// DATATABLE - prefix state name in onchain
const DATATABLE = "Data_"

// Data - struct
type Data struct {
	Key1       string `json:"Key1"`
	Key2       string `json:"Key2"`
	Attribute1 string `json:"Attribute1"`
	Attribute2 string `json:"Attribute2"`
}

/*
 * Init method is called when the Chain code" is instantiated by the blockchain network
 */
func (s *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

/*
 * Invoke method is called as a result of an application request to run the chain code
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	switch function {
	//CreateAdmin
	case "CreateData":
		return createData(stub, args)
	case "UpdateData":
		return updateData(stub, args)
	}
	return shim.Error(fmt.Sprintf("Invoke cannot find function " + function))
}

func createData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key1 := args[0]
	key2 := args[1]
	val1 := args[2]
	val2 := args[3]

	err := util.Createdata(stub, DATATABLE, []string{key1, key2}, &Data{Key1: key1, Key2: key2, Attribute1: val1, Attribute2: val2})
	if err != nil {
		resErr := ResponseError{ResCode: ERR5, Msg: ""}
		return RespondError(resErr)
	}
	resSuc := ResponseSuccess{ResCode: SUCCESS, Msg: ResCodeDict[SUCCESS], Payload: ""}
	return RespondSuccess(resSuc)
}

func updateData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key1 := args[0]
	key2 := args[1]
	val1 := args[2]
	val2 := args[3]

	err := util.UpdateExistingData(stub, DATATABLE, []string{key1, key2}, &Data{Key1: key1, Key2: key2, Attribute1: val1, Attribute2: val2})
	if err != nil {
		resErr := ResponseError{ResCode: ERR5, Msg: ""}
		return RespondError(resErr)
	}
	resSuc := ResponseSuccess{ResCode: SUCCESS, Msg: ResCodeDict[SUCCESS], Payload: ""}
	return RespondSuccess(resSuc)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Chain code
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error creating new Chain code: %s", err)
	}
}
