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
	"github.com/Akachain/akc-go-sdk/hstx/models"
	"github.com/Akachain/akc-go-sdk/util"
	. "github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
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

	// An admin can only create one signed quorum for any given proposal
	quorumResult := new(Quorum)
	queryString := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"^Quorum_\"},\"ProposalID\": \"%s\"}}", ProposalID)
	fmt.Printf("queryString : %s \n", queryString)

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
		fmt.Printf("queryResponse.Value:%v \n", queryResponse.Value)
		_ = json.Unmarshal(queryResponse.Value, quorumResult)
	}

	fmt.Printf("quorumResult:%v \n", quorumResult)

	if 0 == strings.Compare(quorumResult.AdminID, AdminID) {
		resErr := common.ResponseError{common.ERR9, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR9], "Only signed once ", common.GetLine())}
		return common.RespondError(resErr)
	}

	fmt.Printf("Pass if quorum.AdminID == AdminID \n")

	//get data to verify
	rs, errData := Get_data_byid_(stub, ProposalID, models.PROPOSALTABLE)
	dataProposal := rs.(*Proposal)
	fmt.Printf("Pass get data to verify \n")

	if errData != nil {
		resErr := common.ResponseError{common.ERR4, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], errData.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	fmt.Printf("Signature %v\n", Signature)
	fmt.Printf("AdminID %v\n", AdminID)
	fmt.Printf("ProposalID %v\n", ProposalID)
	fmt.Printf("dataProposal %v\n", dataProposal)

	rs, errAd := Get_data_byid_(stub, AdminID, models.ADMINTABLE)

	mapstructure.Decode(rs, admin)
	fmt.Printf("Amdin: %v\n", admin)

	if errAd != nil {
		resErr := common.ResponseError{common.ERR4, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], errAd.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	block, _ := pem.Decode([]byte(admin.PublicKey))
	fmt.Printf("admin.PublicKey %v\n", admin.PublicKey)

	if block == nil {
		resErr := common.ResponseError{common.ERR6, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR6], "block err", common.GetLine())}
		return common.RespondError(resErr)
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		resErr := common.ResponseError{common.ERR6, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR6], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	hashFunc := cp.SHA512
	h := hashFunc.New()
	h.Write([]byte(dataProposal.Data))
	hashed := h.Sum(nil)
	result := rsa.VerifyPKCS1v15(pub, hashFunc, hashed, []byte(Signature))
	fmt.Printf("result %v\n", result)

	if result != nil {
		resErr := common.ResponseError{common.ERR8, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR8], result.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}

	QuorumID := xid.New().String()
	fmt.Printf("QuorumID %v\n", QuorumID)

	err1 := Create_data_(stub, models.QUORUMTABLE, []string{QuorumID}, &Quorum{AdminID: AdminID, QuorumID: QuorumID, ProposalID: ProposalID, Status: "Verify"})
	if err1 != nil {
		resErr := common.ResponseError{common.ERR6, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR6], err1.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], QuorumID}
	return common.RespondSuccess(resSuc)
}

// GetQuorumByID
func (quorum *Quorum) GetQuorumByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	DataID := args[0]
	res := GetDataByID(stub, DataID, quorum, models.QUORUMTABLE)
	return res
}

// GetAllQuorum
func (quorum *Quorum) GetAllQuorum(stub shim.ChaincodeStubInterface) pb.Response {
	res := GetAllData(stub, quorum, models.QUORUMTABLE)
	return res
}

// ------------------- //
