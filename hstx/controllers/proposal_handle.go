package controllers

import (
	"fmt"

	"github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	. "github.com/Akachain/akc-go-sdk/util"
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
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	ProposalID := xid.New().String()
	fmt.Printf("ProposalID %s\n", ProposalID)

	err := Create_data_(stub, models.PROPOSALTABLE, []string{ProposalID}, &Proposal{ProposalID: ProposalID, Data: args[0]})

	if err != nil {
		resErr := common.ResponseError{common.ERR5, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}

	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], ProposalID}
	return common.RespondSuccess(resSuc)
}

// GetProposalByID
func (proposal *Proposal) GetProposalByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	DataID := args[0]
	res := GetDataByID(stub, DataID, proposal, models.PROPOSALTABLE)
	return res
}

// GetAllProposal
func (proposal *Proposal) GetAllProposal(stub shim.ChaincodeStubInterface) pb.Response {
	res := GetAllData(stub, proposal, models.PROPOSALTABLE)
	return res
}

// ------------------- //
