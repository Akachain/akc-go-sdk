package controllers

import (
	cp "crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/Akachain/akc-go-sdk/common"
	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"
)

type Quorum models.Quorum

//High secure transaction Quorum handle
// ------------------- //

func (quorum *Quorum) CreateQuorum(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 3)
	Signature := args[0]
	AdminID := args[1]
	ProposalID := args[2]
	var admin *Admin
	var dataProposal *Proposal

	//check data to verify
	rs, err := util.Getdatabyid(stub, ProposalID, models.PROPOSALTABLE)

	if rs == nil {
		resErr := ResponseError{ERR13, fmt.Sprintf("%s %s %s", ResCodeDict[ERR13], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	mapstructure.Decode(rs, &dataProposal)
	if err != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
		return RespondError(resErr)
	}

	// An admin can only create one signed quorum for any given proposal
	var quorumList = []Quorum{}
	quorumResult := new(Quorum)
	//Select ProposalID from Model Quorum to check exist
	queryString := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"Quorum_\"},\"ProposalID\": \"%s\"}}", ProposalID)
	resultsIterator, err := stub.GetQueryResult(queryString)
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
		err = json.Unmarshal(queryResponse.Value, quorumResult)
		if err != nil {
			//convert JSON eror
			resErr := ResponseError{ERR3, fmt.Sprintf("%s %s %s", ResCodeDict[ERR3], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		quorumList = append(quorumList, *quorumResult)
	}
	// Check AdminID exits in Quorum model
	for _, quorumCompare := range quorumList {
		if 0 == strings.Compare(quorumCompare.AdminID, AdminID) {
			resErr := ResponseError{ERR9, fmt.Sprintf("%s %s %s", ResCodeDict[ERR9], "Only signed once ", GetLine())}
			return RespondError(resErr)
		}
	}
	// Select data from Admin Model
	rs, err = util.Getdatabyid(stub, AdminID, models.ADMINTABLE)

	if rs == nil {
		resErr := ResponseError{ERR14, fmt.Sprintf("%s %s %s", ResCodeDict[ERR14], "", GetLine())}
		return RespondError(resErr)
	}
	mapstructure.Decode(rs, &admin)
	if err != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	//check Status Admin
	if admin.Status != "Active" {
		resErr := ResponseError{ERR15, fmt.Sprintf("%s %s %s", ResCodeDict[ERR15], "", GetLine())}
		return RespondError(resErr)
	}
	block, _ := pem.Decode([]byte(admin.PublicKey))
	if block == nil {
		resErr := ResponseError{ERR6, fmt.Sprintf("%s %s %s", ResCodeDict[ERR6], "block err", GetLine())}
		return RespondError(resErr)
	}
	// Parse to publickey
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		resErr := ResponseError{ERR6, fmt.Sprintf("%s %s %s", ResCodeDict[ERR6], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	// verify signature
	hashFunc := cp.SHA512
	h := hashFunc.New()
	h.Write([]byte(dataProposal.Data))
	hashed := h.Sum(nil)
	err = rsa.VerifyPKCS1v15(pub, hashFunc, hashed, []byte(Signature))

	if err != nil {
		resErr := ResponseError{ERR8, fmt.Sprintf("%s %s %s", ResCodeDict[ERR8], err.Error(), GetLine())}
		return RespondError(resErr)
	}

	QuorumID := stub.GetTxID()
	err = util.Createdata(stub, models.QUORUMTABLE, []string{QuorumID}, &Quorum{AdminID: AdminID, QuorumID: QuorumID, ProposalID: ProposalID, Status: "Verify"})

	if err != nil {
		resErr := ResponseError{ERR5, fmt.Sprintf("%s %s %s", ResCodeDict[ERR5], err.Error(), GetLine())}
		return RespondError(resErr)
	}

	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], QuorumID}
	return RespondSuccess(resSuc)
}

//CreateReject
func (quorum *Quorum) CreateReject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 3)
	Signature := args[0]
	AdminID := args[1]
	ProposalID := args[2]
	var admin *Admin
	var dataProposal *Proposal

	//check data to verify
	rs, err := util.Getdatabyid(stub, ProposalID, models.PROPOSALTABLE)

	if rs == nil {
		resErr := ResponseError{ERR13, fmt.Sprintf("%s %s %s", ResCodeDict[ERR13], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	mapstructure.Decode(rs, &dataProposal)
	if err != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
		return RespondError(resErr)
	}

	// An admin can only create one signed quorum for any given proposal
	var quorumList = []Quorum{}
	quorumResult := new(Quorum)
	//Select ProposalID from Model Quorum to check exist
	queryString := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"Quorum_\"},\"ProposalID\": \"%s\"}}", ProposalID)
	resultsIterator, err := stub.GetQueryResult(queryString)
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
		err = json.Unmarshal(queryResponse.Value, quorumResult)
		if err != nil {
			//convert JSON eror
			resErr := ResponseError{ERR3, fmt.Sprintf("%s %s %s", ResCodeDict[ERR3], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		quorumList = append(quorumList, *quorumResult)
	}
	// Check AdminID exits in Quorum model
	for _, quorumCompare := range quorumList {
		if 0 == strings.Compare(quorumCompare.AdminID, AdminID) && quorumCompare.Status == "Verify" {
			resErr := ResponseError{ERR17, fmt.Sprintf("%s %s %s", ResCodeDict[ERR17], "", GetLine())}
			return RespondError(resErr)
		}
	}

	// Check AdminID exits in Quorum model
	for _, quorumCompare := range quorumList {
		if 0 == strings.Compare(quorumCompare.AdminID, AdminID) && quorumCompare.Status == "Reject" {
			resErr := ResponseError{ERR18, fmt.Sprintf("%s %s %s", ResCodeDict[ERR18], "", GetLine())}
			return RespondError(resErr)
		}
	}

	// Select data from Admin Model
	rsadmin, err := util.Getdatabyid(stub, AdminID, models.ADMINTABLE)

	if rsadmin == nil {
		resErr := ResponseError{ERR14, fmt.Sprintf("%s %s %s", ResCodeDict[ERR14], "", GetLine())}
		return RespondError(resErr)
	}
	mapstructure.Decode(rsadmin, &admin)
	if err != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	//check Status Admin
	if admin.Status != "Active" {
		resErr := ResponseError{ERR15, fmt.Sprintf("%s %s %s", ResCodeDict[ERR15], "", GetLine())}
		return RespondError(resErr)
	}
	block, _ := pem.Decode([]byte(admin.PublicKey))
	if block == nil {
		resErr := ResponseError{ERR6, fmt.Sprintf("%s %s %s", ResCodeDict[ERR6], "block err", GetLine())}
		return RespondError(resErr)
	}
	// Parse to publickey
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		resErr := ResponseError{ERR6, fmt.Sprintf("%s %s %s", ResCodeDict[ERR6], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	// verify signature
	hashFunc := cp.SHA512
	h := hashFunc.New()
	h.Write([]byte(dataProposal.Data))
	hashed := h.Sum(nil)
	err = rsa.VerifyPKCS1v15(pub, hashFunc, hashed, []byte(Signature))

	if err != nil {
		resErr := ResponseError{ERR8, fmt.Sprintf("%s %s %s", ResCodeDict[ERR8], err.Error(), GetLine())}
		return RespondError(resErr)
	}

	QuorumID := stub.GetTxID()
	err = util.Createdata(stub, models.QUORUMTABLE, []string{QuorumID}, &Quorum{AdminID: AdminID, QuorumID: QuorumID, ProposalID: ProposalID, Status: "Reject"})
	if err != nil {
		resErr := ResponseError{ERR5, fmt.Sprintf("%s %s %s", ResCodeDict[ERR5], err.Error(), GetLine())}
		return RespondError(resErr)
	}

	err = UpdateProposal(stub, ProposalID, "Reject")
	fmt.Printf("err: %v\n", err)
	if err != nil {
		resErr := ResponseError{ERR5, fmt.Sprintf("%s %s %s", ResCodeDict[ERR5], err.Error(), GetLine())}
		return RespondError(resErr)
	}

	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], QuorumID}
	return RespondSuccess(resSuc)
}

// GetQuorumByID
func (quorum *Quorum) GetQuorumByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 1)

	DataID := args[0]
	res := util.GetDataByID(stub, DataID, quorum, models.QUORUMTABLE)
	return res
}

// GetAllQuorum
func (quorum *Quorum) GetAllQuorum(stub shim.ChaincodeStubInterface) pb.Response {
	res := util.GetAllData(stub, quorum, models.QUORUMTABLE)
	return res
}

//GetQuorumByProposalID
func (quorum *Quorum) GetQuorumByProposalID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	ProposalID := args[0]
	var quorumList = []Quorum{}
	quorumResult := new(Quorum)
	queryString := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"Quorum_\"},\"ProposalID\": \"%s\"}}", ProposalID)
	resultsIterator, err := stub.GetQueryResult(queryString)
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
		err = json.Unmarshal(queryResponse.Value, quorumResult)
		if err != nil {
			//convert JSON eror
			resErr := ResponseError{ERR3, fmt.Sprintf("%s %s %s", ResCodeDict[ERR3], err.Error(), GetLine())}
			return RespondError(resErr)
		}
		quorumList = append(quorumList, *quorumResult)
	}
	quorumJson, err2 := json.Marshal(quorumList)
	if err2 != nil {
		//convert JSON eror
		resErr := common.ResponseError{common.ERR3, common.ResCodeDict[common.ERR3]}
		return common.RespondError(resErr)
	}
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(quorumJson)}
	return common.RespondSuccess(resSuc)
}

// ------------------- //
