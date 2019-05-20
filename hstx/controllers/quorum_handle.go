package controllers

import (
	cp "crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
)

type Quorum models.Quorum

//High secure transaction Quorum handle
// ------------------- //

func (quorum *Quorum) CreateQuorum(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		//Invalid arguments
		resErr := ResponseError{ERR2, ResCodeDict[ERR2]}
		return RespondError(resErr)
	}
	Signature := args[0]
	AdminID := args[1]
	ProposalID := args[2]
	var admin *Admin

	//check Only signed once
	quorumResult := new(Quorum)
	queryString := fmt.Sprintf("{\"selector\": {\"_id\": {\"$regex\": \"^Quorum_\"},\"ProposalID\": \"%s\"}}", ProposalID)
	fmt.Printf("queryString : %s \n", queryString)

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
		fmt.Printf("queryResponse.Value:%v \n", queryResponse.Value)
		_ = json.Unmarshal(queryResponse.Value, quorumResult)
	}

	fmt.Printf("quorumResult:%v \n", quorumResult)

	if 0 == strings.Compare(quorumResult.AdminID, AdminID) {
		resErr := ResponseError{ERR9, fmt.Sprintf("%s %s %s", ResCodeDict[ERR9], "Only signed once ", GetLine())}
		return RespondError(resErr)
	}

	fmt.Printf("Pass if quorum.AdminID == AdminID \n")

	//get data to verify
	rs, errData := util.Getdatabyid(stub, ProposalID, models.PROPOSALTABLE)
	dataProposal := rs.(*Proposal)
	fmt.Printf("Pass get data to verify \n")

	if errData != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], errData.Error(), GetLine())}
		return RespondError(resErr)
	}
	fmt.Printf("Signature %v\n", Signature)
	fmt.Printf("AdminID %v\n", AdminID)
	fmt.Printf("ProposalID %v\n", ProposalID)
	fmt.Printf("dataProposal %v\n", dataProposal)

	rs, errAd := util.Getdatabyid(stub, AdminID, models.ADMINTABLE)

	mapstructure.Decode(rs, admin)
	fmt.Printf("Amdin: %v\n", admin)

	if errAd != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s %s", ResCodeDict[ERR4], errAd.Error(), GetLine())}
		return RespondError(resErr)
	}
	block, _ := pem.Decode([]byte(admin.PublicKey))
	fmt.Printf("admin.PublicKey %v\n", admin.PublicKey)

	if block == nil {
		resErr := ResponseError{ERR6, fmt.Sprintf("%s %s %s", ResCodeDict[ERR6], "block err", GetLine())}
		return RespondError(resErr)
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		resErr := ResponseError{ERR6, fmt.Sprintf("%s %s %s", ResCodeDict[ERR6], err.Error(), GetLine())}
		return RespondError(resErr)
	}
	hashFunc := cp.SHA512
	h := hashFunc.New()
	h.Write([]byte(dataProposal.Data))
	hashed := h.Sum(nil)
	result := rsa.VerifyPKCS1v15(pub, hashFunc, hashed, []byte(Signature))
	fmt.Printf("result %v\n", result)

	if result != nil {
		resErr := ResponseError{ERR8, fmt.Sprintf("%s %s %s", ResCodeDict[ERR8], result.Error(), GetLine())}
		return RespondError(resErr)
	}

	QuorumID := xid.New().String()
	fmt.Printf("QuorumID %v\n", QuorumID)

	err1 := util.Createdata(stub, models.QUORUMTABLE, []string{QuorumID}, &Quorum{AdminID: AdminID, QuorumID: QuorumID, ProposalID: ProposalID, Status: "Verify"})
	if err1 != nil {
		resErr := ResponseError{ERR6, fmt.Sprintf("%s %s %s", ResCodeDict[ERR6], err1.Error(), GetLine())}
		return RespondError(resErr)
	}
	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], QuorumID}
	return RespondSuccess(resSuc)
}

// GetQuorumByID
func (quorum *Quorum) GetQuorumByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		//Invalid arguments
		resErr := ResponseError{ERR2, ResCodeDict[ERR2]}
		return RespondError(resErr)
	}
	DataID := args[0]
	res := util.GetDataByID(stub, DataID, quorum, models.QUORUMTABLE)
	return res
}

// GetAllQuorum
func (quorum *Quorum) GetAllQuorum(stub shim.ChaincodeStubInterface) pb.Response {
	res := util.GetAllData(stub, quorum, models.QUORUMTABLE)
	return res
}

// ------------------- //
