package controllers

import (
	"fmt"

	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Proposal models.Proposal

//High secure transaction Proposal handle
// ------------------- //

//Create Proposal
func (proposal *Proposal) CreateProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 1)
	ProposalID := stub.GetTxID()
	err := util.Createdata(stub, models.PROPOSALTABLE, []string{ProposalID}, &Proposal{ProposalID: ProposalID, Data: args[0]})
	if err != nil {
		resErr := ResponseError{ERR5, fmt.Sprintf("%s %s %s", ResCodeDict[ERR5], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], ProposalID}
	return RespondSuccess(resSuc)
}

// GetProposalByID
func (proposal *Proposal) GetProposalByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 1)
	DataID := args[0]
	res := util.GetDataByID(stub, DataID, proposal, models.PROPOSALTABLE)
	return res
}

// GetAllProposal
func (proposal *Proposal) GetAllProposal(stub shim.ChaincodeStubInterface) pb.Response {
	res := util.GetAllData(stub, proposal, models.PROPOSALTABLE)
	return res
}

// ------------------- //
