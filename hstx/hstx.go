package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	. "github.com/Akachain/akc-go-sdk/common"
	ctl "github.com/Akachain/akc-go-sdk/hstx/controllers"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode implementation
type Chaincode struct {
}

/*
 * The Init method is called when the Chain code" is instantiated by the blockchain network
 */
func (s *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// The invokeFunction returns
	Admin1 := "Admin1"
	pubKey1, _ := ioutil.ReadFile("./sample/pk1.pem")
	pk1 := base64.StdEncoding.EncodeToString(pubKey1)

	Admin2 := "Admin2"
	pubKey2, _ := ioutil.ReadFile("./sample/pk2.pem")
	pk2 := base64.StdEncoding.EncodeToString(pubKey2)

	Admin3 := "Admin3"
	pubKey3, _ := ioutil.ReadFile("./sample/pk3.pem")
	pk3 := base64.StdEncoding.EncodeToString(pubKey3)

	rs1 := controller_admin.CreateAdmin(stub, []string{Admin1, pk1})
	rs2 := controller_admin.CreateAdmin(stub, []string{Admin2, pk2})
	rs3 := controller_admin.CreateAdmin(stub, []string{Admin3, pk3})

	if rs1.Status != shim.OK || rs2.Status != shim.OK || rs3.Status != shim.OK {
		return shim.Error("Init chaincode with 3 admin fail")
	}
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
	//CreateAdmin
	case "CreateAdmin":
		return controller_admin.CreateAdmin(stub, args)
	//UpdateAdmin -> For testing
	case "UpdateAdmin":
		return controller_admin.UpdateAdmin(stub, args)
	//CreateProposal
	case "CreateProposal":
		return controller_proposal.CreateProposal(stub, args)
	//CreateQuorum
	case "CreateQuorum":
		return controller_quorum.CreateQuorum(stub, args)
	//CreateCommit
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
	// GetAdminByID
	case "GetAdminByID":
		return controller_admin.GetAdminByID(stub, args)
	// GetAllAdmin
	case "GetAllAdmin":
		return controller_admin.GetAllAdmin(stub)
	// GetProposalByID
	case "GetProposalByID":
		return controller_proposal.GetProposalByID(stub, args)
	// GetAllProposal
	case "GetAllProposal":
		return controller_proposal.GetAllProposal(stub)
	// GetQuorumByID
	case "GetQuorumByID":
		return controller_quorum.GetQuorumByID(stub, args)
	// GetAllQuorum
	case "GetAllQuorum":
		return controller_quorum.GetAllQuorum(stub)
	// GetCommitByID
	case "GetCommitByID":
		return controller_commit.GetCommitByID(stub, args)
	// GetAllCommit
	case "GetAllCommit":
		return controller_commit.GetAllCommit(stub)
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
