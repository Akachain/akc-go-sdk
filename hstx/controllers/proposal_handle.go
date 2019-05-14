package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
	"gitlab.com/akachain/akc-go-sdk/common"
	"gitlab.com/akachain/akc-go-sdk/hstx/models"
)

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

	err := create_data_(stub, models.PROPOSALTABLE, []string{ProposalID}, &Proposal{ProposalID: ProposalID, Data: args[0]})

	if err != nil {
		resErr := common.ResponseError{common.ERR5, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}

	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], ProposalID}
	return common.RespondSuccess(resSuc)
}

func (proposal *Proposal) GetProposalByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	ProposalID := args[0]

	rs, err := get_data_byid_(stub, ProposalID, models.PROPOSALTABLE)

	mapstructure.Decode(rs, proposal)
	fmt.Printf("Proposal: %v\n", proposal)

	bytes, err := json.Marshal(proposal)
	if err != nil {
		//Convert Json Fail
		resErr := common.ResponseError{common.ERR3, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	fmt.Printf("Response: %s\n", string(bytes))

	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(bytes)}
	return common.RespondSuccess(resSuc)
}

func (proposal *Proposal) GetAllProposal(stub shim.ChaincodeStubInterface) pb.Response {
	proposalbytes, err := get_all_data_(stub, models.PROPOSALTABLE)

	proposal = new(Proposal)
	Proposallist := []*Proposal{}

	for row_json_bytes := range proposalbytes {
		proposal = new(Proposal)
		err = json.Unmarshal(row_json_bytes, proposal)
		if err != nil {

			resErr := common.ResponseError{common.ERR6, common.ResCodeDict[common.ERR6]}
			return common.RespondError(resErr)
		}
		Proposallist = append(Proposallist, proposal)
	}

	if err != nil {
		//Get data eror
		resErr := common.ResponseError{common.ERR3, common.ResCodeDict[common.ERR3]}
		return common.RespondError(resErr)
	}
	proposalJson, err2 := json.Marshal(Proposallist)
	if err2 != nil {
		//convert JSON eror
		resErr := common.ResponseError{common.ERR6, common.ResCodeDict[common.ERR6]}
		return common.RespondError(resErr)
	}

	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(proposalJson)}
	return common.RespondSuccess(resSuc)
}

// ------------------- //
