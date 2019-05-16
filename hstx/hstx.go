package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	. "gitlab.com/akachain/akc-go-sdk/common"
	ctl "gitlab.com/akachain/akc-go-sdk/hstx/controllers"
)

// Chaincode implementation
type Chaincode struct {
}

/*
 * The Init method is called when the Chain code" is instantiated by the blockchain network
 */
func (s *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

var controller_admin ctl.Admin
var controller_proposal ctl.Proposal
var controller_quorum ctl.Quorum
var controller_commit ctl.Commit

/*
 * The Invoke method is called as a result of an application request to run the chain code
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	Logger.Info("########### Hstx Invoke ###########")
	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	switch function {
	//Create CreateAdmin
	case "CreateAdmin":
		return controller_admin.CreateAdmin(stub, args)
	//Create CreateProposal
	case "CreateProposal":
		return controller_proposal.CreateProposal(stub, args)
		//Create CreateQuorum
	case "CreateQuorum":
		return controller_quorum.CreateQuorum(stub, args)
		//Create CreateCommit
	case "CreateCommit":
		return controller_commit.CreateCommit(stub, args)
	default:
		return t.Query(stub)
	}
}

// Query callback representing the query of a chaincode
func (t *Chaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	Logger.Info("########### Hstx Query ###########")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	// GetDataID
	case "GetDataByID":
		return ctl.GetDataByID(stub, args)
		// GetAllData
	case "GetAllData":
		return ctl.GetAllData(stub, args)
	}

	return shim.Error(fmt.Sprintf("[Hstx Chaincode] Invoke and Query not find function " + function))
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Chain code
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error creating new Chain code: %s", err)
	}
}
