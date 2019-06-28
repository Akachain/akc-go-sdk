package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Akachain/akc-go-sdk/common"
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
	Data := args[0]
	txTimeStamp, _ := stub.GetTxTimestamp()
	CreateDatetime := time.Unix(txTimeStamp.Seconds, int64(txTimeStamp.Nanos)).UTC()

	err := util.Createdata(stub, models.PROPOSALTABLE, []string{ProposalID}, &Proposal{ProposalID: ProposalID, Data: Data, CreateDateTime: CreateDatetime.String()})
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

//GetProposalNotSign
func (proposal *Proposal) GetProposalNotSign(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	AdminID := args[0]
	var proposalList = []Proposal{}
	var proposalListResponse []string
	proposalResult := new(Proposal)

	queryStringProposal := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"Proposal_\"}}}")

	resultsIterator, errProposal := stub.GetQueryResult(queryStringProposal)
	if errProposal != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], errProposal.Error(), GetLine())}
		return RespondError(resErr)
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		err = json.Unmarshal(queryResponse.Value, proposalResult)
		if err != nil {
			//convert JSON eror
			resErr := ResponseError{ERR3, fmt.Sprintf("%s %s %s", ResCodeDict[ERR3], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		proposalList = append(proposalList, *proposalResult)
	}

	for _, proposal := range proposalList {
		var quorumList = []Quorum{}
		quorumResult := new(Quorum)
		queryString := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"Quorum_\"},\"ProposalID\": \"%s\",\"AdminID\": \"%s\"}}", proposal.ProposalID, AdminID)
		resultsIterator, errQuorum := stub.GetQueryResult(queryString)
		if errQuorum != nil {
			resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], errQuorum.Error(), GetLine())}
			return RespondError(resErr)
		}
		defer resultsIterator.Close()

		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
				return RespondError(resErr)
			}
			err = json.Unmarshal(queryResponse.Value, quorumResult)
			if err != nil {
				//convert JSON eror
				resErr := ResponseError{ERR3, fmt.Sprintf("%s %s %s", ResCodeDict[ERR3], err.Error(), GetLine())}
				return RespondError(resErr)
			}
			quorumList = append(quorumList, *quorumResult)
		}

		// Check ProposalID exits in Quorum model
		if len(quorumList) == 0 {
			proposalListResponse = append(proposalListResponse, proposal.ProposalID)
		}
	}
	proposalJson, err2 := json.Marshal(proposalListResponse)
	if err2 != nil {
		//convert JSON eror
		resErr := common.ResponseError{common.ERR3, common.ResCodeDict[common.ERR3]}
		return common.RespondError(resErr)
	}
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(proposalJson)}
	return common.RespondSuccess(resSuc)
}

// ------------------- //
