package controllers

import (
	"encoding/json"
	"fmt"

	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Commit models.Commit

//High secure transaction Commit handle
// ------------------- //
//Create Commit
func (commit *Commit) CreateCommit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 2)
	AdminID := args[0]
	ProposalID := args[1]
	var admin *Admin
	admin = new(Admin)
	Adminlist := []*Admin{}
	adminbytes, err := util.Getalldata(stub, models.ADMINTABLE)
	if err != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	for row_json_bytes := range adminbytes {
		admin = new(Admin)
		err = json.Unmarshal(row_json_bytes, admin)
		//convert JSON eror
		if err != nil {
			resErr := ResponseError{ERR3, fmt.Sprintf("%s %s %s", ResCodeDict[ERR3], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		Adminlist = append(Adminlist, admin)
	}

	var quorumList = []Quorum{}
	quorumResutl := new(Quorum)
	commitResutl := new(Commit)
	//check ProposalID exist in Quorum
	queryStringQuorum := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"Quorum_\"},\"ProposalID\": \"%s\"}}", ProposalID)

	Logger.Debug("queryStringQuorum: %v", queryStringQuorum)
	resultsIterator, err := stub.GetQueryResult(queryStringQuorum)
	if err != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		err = json.Unmarshal(queryResponse.Value, quorumResutl)
		if err != nil {
			//convert JSON eror
			resErr := ResponseError{ERR3, fmt.Sprintf("%s %s %s", ResCodeDict[ERR3], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		quorumList = append(quorumList, *quorumResutl)
	}
	if quorumResutl.ProposalID == "" {
		resErr := ResponseError{ERR12, fmt.Sprintf("%s %s %s", ResCodeDict[ERR12], "", GetLine())}
		return RespondError(resErr)
	}

	//check Only Commit once
	queryStringCommit := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"Commit_\"},\"ProposalID\": \"%s\"}}", ProposalID)
	Logger.Debug("queryStringCommit: %v", queryStringCommit)

	resultsIterator, err = stub.GetQueryResult(queryStringCommit)
	if err != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		err = json.Unmarshal(queryResponse.Value, commitResutl)
		if err != nil {
			//convert JSON eror
			resErr := ResponseError{ERR3, fmt.Sprintf("%s %s %s", ResCodeDict[ERR3], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		quorumList = append(quorumList, *quorumResutl)
	}

	if commitResutl.CommitID != "" && commitResutl.Status == "Verify" {
		resErr := ResponseError{ERR11, fmt.Sprintf("%s %s %s", ResCodeDict[ERR11], "", GetLine())}
		return RespondError(resErr)
	}

	count := 0
	var quorumIDList []string
	for _, quorum := range quorumList {
		if quorum.Status == "Verify" {
			for _, admin := range Adminlist {
				if quorum.AdminID == admin.AdminID && admin.Status == "Active" {
					count = count + 1
					quorumIDList = append(quorumIDList, quorum.QuorumID)
				}
			}
		} else if quorum.Status == "Reject" {
			Logger.Debug("Reject: %v", count)
			AdminReject := fmt.Sprintf("Proposal: %s  was rejected by AdminID: %s", quorum.ProposalID, quorum.AdminID)
			resErr := ResponseError{ERR10, fmt.Sprintf("%s %s %s", ResCodeDict[ERR16], AdminReject, GetLine())}
			return RespondError(resErr)
		}
	}

	if count < 3 || len(quorumIDList) < 3 {
		Logger.Debug("Not Enough quorum: %v", count)
		resErr := ResponseError{ERR10, fmt.Sprintf("%s %s %s", ResCodeDict[ERR10], "[]", GetLine())}
		return RespondError(resErr)
	}
	CommitID := stub.GetTxID()
	Logger.Debug("CommitID Return: %v", CommitID)

	err = util.Createdata(stub, models.COMMITTABLE, []string{CommitID}, &Commit{CommitID: string(CommitID), AdminID: AdminID, ProposalID: ProposalID, QuorumList: quorumIDList, Status: "Verify"})
	if err != nil {
		resErr := ResponseError{ERR5, fmt.Sprintf("%s %s %s", ResCodeDict[ERR5], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	err = UpdateProposal(stub, ProposalID, "Approve")
	fmt.Printf("err: %v\n", err)
	if err != nil {
		resErr := ResponseError{ERR5, fmt.Sprintf("%s %s %s", ResCodeDict[ERR5], err.Error(), GetLine())}
		return RespondError(resErr)
	}

	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], CommitID}
	return RespondSuccess(resSuc)
}

// GetCommitByID
func (commit *Commit) GetCommitByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 1)

	DataID := args[0]
	res := util.GetDataByID(stub, DataID, commit, models.COMMITTABLE)
	return res
}

// GetAllCommit
func (commit *Commit) GetAllCommit(stub shim.ChaincodeStubInterface) pb.Response {
	res := util.GetAllData(stub, commit, models.COMMITTABLE)
	return res
}

// ------------------- //
