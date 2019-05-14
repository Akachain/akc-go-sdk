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

//High secure transaction Commit handle
// ------------------- //
//Create Commit
func (commit *Commit) CreateCommit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	ProposalID := args[0]
	var admin *Admin
	adminbytes, err := get_all_data_(stub, models.ADMINTABLE)

	admin = new(Admin)
	Adminlist := []*Admin{}

	for row_json_bytes := range adminbytes {
		admin = new(Admin)
		err = json.Unmarshal(row_json_bytes, admin)
		if err != nil {

			resErr := common.ResponseError{common.ERR6, common.ResCodeDict[common.ERR6]}
			return common.RespondError(resErr)
		}
		Adminlist = append(Adminlist, admin)
	}

	var quorumList = []Quorum{}
	quorumResutl := new(Quorum)
	commitResutl := new(Commit)
	queryString := fmt.Sprintf("{\"selector\": {\"ProposalID\": \"%s\"}}", ProposalID)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		resErr := common.ResponseError{common.ERR4, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			resErr := common.ResponseError{common.ERR4, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())}
			return common.RespondError(resErr)
		}
		_ = json.Unmarshal(queryResponse.Value, commitResutl)
		_ = json.Unmarshal(queryResponse.Value, quorumResutl)
		quorumList = append(quorumList, *quorumResutl)
	}
	if quorumResutl.ProposalID == "" {
		resErr := common.ResponseError{common.ERR12, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR12], "", common.GetLine())}
		return common.RespondError(resErr)
	}
	if commitResutl.CommitID != "" && commitResutl.Status == "Verify" {
		resErr := common.ResponseError{common.ERR11, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR11], "", common.GetLine())}
		return common.RespondError(resErr)
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
		resErr := common.ResponseError{common.ERR10, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR10], "[]", common.GetLine())}
		return common.RespondError(resErr)
	}
	CommitID := xid.New().String()
	fmt.Printf("CommitID %v\n", CommitID)

	err1 := create_data_(stub, models.COMMITTABLE, []string{CommitID}, &Commit{CommitID: string(CommitID), ProposalID: ProposalID, QuorumID: quorumIDList, Status: "Verify"})
	if err1 != nil {
		resErr := common.ResponseError{common.ERR6, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR6], err1.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], CommitID}
	return common.RespondSuccess(resSuc)
}

//GetCommitByID
func (commit *Commit) GetCommitByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	CommitID := args[0]

	rs, err := get_data_byid_(stub, CommitID, models.COMMITTABLE)

	mapstructure.Decode(rs, commit)
	fmt.Printf("Commit: %v\n", commit)

	bytes, err := json.Marshal(commit)
	if err != nil {
		//Convert Json Fail
		resErr := common.ResponseError{common.ERR3, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	fmt.Printf("Response: %s\n", string(bytes))

	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(bytes)}
	return common.RespondSuccess(resSuc)
}

//GetAllCommit
func (commit *Commit) GetAllCommit(stub shim.ChaincodeStubInterface) pb.Response {
	commitbytes, err := get_all_data_(stub, models.COMMITTABLE)

	commit = new(Commit)
	Commitlist := []*Commit{}

	for row_json_bytes := range commitbytes {
		commit = new(Commit)
		err = json.Unmarshal(row_json_bytes, commit)
		if err != nil {

			resErr := common.ResponseError{common.ERR6, common.ResCodeDict[common.ERR6]}
			return common.RespondError(resErr)
		}
		Commitlist = append(Commitlist, commit)
	}

	commitJson, err2 := json.Marshal(Commitlist)
	if err2 != nil {
		//convert JSON eror
		resErr := common.ResponseError{common.ERR6, common.ResCodeDict[common.ERR6]}
		return common.RespondError(resErr)
	}

	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(commitJson)}
	return common.RespondSuccess(resSuc)
}

// ------------------- //
