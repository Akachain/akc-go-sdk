package controllers

import (
	"encoding/json"
	"fmt"

	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/rs/xid"
)

type Commit models.Commit

//High secure transaction Commit handle
// ------------------- //
//Create Commit
func (commit *Commit) CreateCommit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		//Invalid arguments
		resErr := ResponseError{ERR2, ResCodeDict[ERR2]}
		return RespondError(resErr)
	}
	ProposalID := args[0]
	var admin *Admin
	adminbytes, err := util.Getalldata(stub, models.ADMINTABLE)

	admin = new(Admin)
	Adminlist := []*Admin{}

	for row_json_bytes := range adminbytes {
		admin = new(Admin)
		err = json.Unmarshal(row_json_bytes, admin)
		if err != nil {

			resErr := ResponseError{ERR6, ResCodeDict[ERR6]}
			return RespondError(resErr)
		}
		Adminlist = append(Adminlist, admin)
	}

	var quorumList = []Quorum{}
	quorumResutl := new(Quorum)
	commitResutl := new(Commit)

	//check ProposalID exist in Quorum
	queryStringQuorum := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"^Quorum_\"},\"ProposalID\": \"%s\"}}", ProposalID)
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
		errQuo := json.Unmarshal(queryResponse.Value, quorumResutl)
		if errQuo != nil {
			//convert JSON eror
			resErr := ResponseError{ERR6, fmt.Sprintf("%s %s %s", ResCodeDict[ERR6], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		quorumList = append(quorumList, *quorumResutl)
	}
	if quorumResutl.ProposalID == "" {
		resErr := ResponseError{ERR12, fmt.Sprintf("%s %s %s", ResCodeDict[ERR12], "", GetLine())}
		return RespondError(resErr)
	}

	//check Only Commit once
	queryStringCommit := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"^Commit_\"},\"ProposalID\": \"%s\"}}", ProposalID)
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
		errCommit := json.Unmarshal(queryResponse.Value, commitResutl)
		if errCommit != nil {
			//convert JSON eror
			resErr := ResponseError{ERR6, fmt.Sprintf("%s %s %s", ResCodeDict[ERR6], err.Error(), GetLine())}
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
				if quorum.AdminID == admin.AdminID {
					count = count + 1
					quorumIDList = append(quorumIDList, quorum.QuorumID)
				}
			}
		}
	}
	fmt.Printf("count:%v \n", count)
	fmt.Printf("quorumIDList:%v \n", quorumIDList)
	fmt.Printf("adminList:%v \n", Adminlist)

	if count < 3 {
		fmt.Printf("Not Enough quorum \n")
		resErr := ResponseError{ERR10, fmt.Sprintf("%s %s %s", ResCodeDict[ERR10], "[]", GetLine())}
		return RespondError(resErr)
	}
	CommitID := xid.New().String()
	fmt.Printf("CommitID %v\n", CommitID)

	err1 := util.Createdata(stub, models.COMMITTABLE, []string{CommitID}, &Commit{CommitID: string(CommitID), ProposalID: ProposalID, QuorumID: quorumIDList, Status: "Verify"})
	if err1 != nil {
		resErr := ResponseError{ERR6, fmt.Sprintf("%s %s %s", ResCodeDict[ERR6], err1.Error(), GetLine())}
		return RespondError(resErr)
	}
	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], CommitID}
	return RespondSuccess(resSuc)
}

// GetCommitByID
func (commit *Commit) GetCommitByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		//Invalid arguments
		resErr := ResponseError{ERR2, ResCodeDict[ERR2]}
		return RespondError(resErr)
	}
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
