package controllers

import (
	"fmt"

	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/rs/xid"
)

type Proposal models.Proposal

//High secure transaction Proposal handle
// ------------------- //

//Create Proposal
func (proposal *Proposal) CreateProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		//Invalid arguments
		resErr := ResponseError{ERR2, ResCodeDict[ERR2]}
		return RespondError(resErr)
	}
	ProposalID := xid.New().String()
	fmt.Printf("ProposalID %s\n", ProposalID)

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
	if len(args) != 1 {
		//Invalid arguments
		resErr := ResponseError{ERR2, ResCodeDict[ERR2]}
		return RespondError(resErr)
	}
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
