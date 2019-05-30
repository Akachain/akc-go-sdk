package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

//test init chaincode with 3 admin
func TestInit(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
	//pubKey1, _ := ioutil.ReadFile("./sample/pk1.pem")
	//pk1 := base64.StdEncoding.EncodeToString(pubKey1)
	// pubKey2, _ := ioutil.ReadFile("./sample/pk2.pem")
	// pk2 := base64.StdEncoding.EncodeToString(pubKey2)
	// pubKey3, _ := ioutil.ReadFile("./sample/pk2.pem")
	// pk3 := base64.StdEncoding.EncodeToString(pubKey3)

	rs := util.MockInitTransaction(t, stub, [][]byte{[]byte("")})
	assert.Equal(t, "", rs)

	// rs1 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("GetAllAdmin")})
	// fmt.Println(rs1)
	// var r1 InvokeResponse
	// err := json.Unmarshal([]byte(rs1), &r1)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var r2 models.Admin
	// json.Unmarshal([]byte(r1.Msg), &r2)
	// assert.Equal(t, pk1, r2.PublicKey)
}
func TestAdmin(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
	adminName := "Admin1"
	pubKey, _ := ioutil.ReadFile("./sample/pk1.pem")
	pk := base64.StdEncoding.EncodeToString(pubKey)

	// Create a new Admin - automatically fail if not succeess
	fmt.Println("Invoke CreateAdmin ", adminName)
	rs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin1"), []byte(pk)})
	// The invokeFunction returns adminID key
	var r InvokeResponse
	json.Unmarshal([]byte(rs), &r)

	// Check if the created admin exist in the ledger
	compositeKey, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r.Rows})
	state, _ := stub.GetState(compositeKey)
	var ad models.Admin
	json.Unmarshal([]byte(state), &ad)

	// Check if the created admin information is correct
	assert.Equal(t, adminName, ad.Name)
	assert.Equal(t, pk, ad.PublicKey)

}
func TestProposal(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
	prData := "Update Money"

	// Create a new Proposal - automatically fail if not succeess
	rs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateProposal"), []byte("Update Money")})

	// The invokeFunction returns ProposalID key
	var r InvokeResponse
	json.Unmarshal([]byte(rs), &r)

	// Check if the created proposal exist in the ledger
	compositeKey, _ := stub.CreateCompositeKey(models.PROPOSALTABLE, []string{r.Rows})
	state, _ := stub.GetState(compositeKey)
	var pr models.Proposal
	json.Unmarshal([]byte(state), &pr)

	// Check if the created proposal information is correct
	assert.Equal(t, prData, pr.Data)
}

func TestQuorum(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
	pubKey1, _ := ioutil.ReadFile("./sample/pk1.pem")
	signature1, _ := ioutil.ReadFile("./sample/signature1.txt")

	// Create a new Admin - automatically fail if not succeess
	rs1 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin1"), []byte(pubKey1)})
	// The invokeFunction returns adminID key
	var r1 InvokeResponse
	json.Unmarshal([]byte(rs1), &r1)
	// Check if the created admin exist in the ledger
	compositeKey1, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r1.Rows})
	state1, _ := stub.GetState(compositeKey1)
	var ad models.Admin
	json.Unmarshal([]byte(state1), &ad)

	// Create a new Proposal - automatically fail if not succeess
	rs2 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateProposal"), []byte("Update Money")})
	// The invokeFunction returns ProposalID key
	var r2 InvokeResponse
	json.Unmarshal([]byte(rs2), &r2)
	// Check if the created Proposal exist in the ledger
	compositeKey2, _ := stub.CreateCompositeKey(models.PROPOSALTABLE, []string{r2.Rows})
	state2, _ := stub.GetState(compositeKey2)
	var pr models.Proposal
	json.Unmarshal([]byte(state2), &pr)

	rs3 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature1, []byte(ad.AdminID), []byte(pr.ProposalID)})
	var r3 InvokeResponse
	json.Unmarshal([]byte(rs3), &r3)
	compositeKey3, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r3.Rows})
	state3, _ := stub.GetState(compositeKey3)
	var qr models.Quorum
	json.Unmarshal([]byte(state3), &qr)

	// confirm create quorum success
	assert.Equal(t, "200", r3.Status)
	// confirm data in quorum struct
	assert.Equal(t, pr.ProposalID, qr.ProposalID)
	assert.Equal(t, "Verify", qr.Status)
	assert.Equal(t, ad.AdminID, qr.AdminID)

}

func TestQuorum_Err(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
	pubKey1, _ := ioutil.ReadFile("./sample/pk1.pem")
	signature1, _ := ioutil.ReadFile("./sample/signature1.txt")
	signatureFail, _ := ioutil.ReadFile("./sample/signatureFail.txt")
	ProposalIDFail := "ProposalID no thing"
	AdminIDFail := "AdminID no thing"

	// Create a new Admin - automatically fail if not succeess
	rs1 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin1"), []byte(pubKey1)})
	// The invokeFunction returns adminID key
	var r1 InvokeResponse
	json.Unmarshal([]byte(rs1), &r1)
	// Check if the created admin exist in the ledger
	compositeKey1, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r1.Rows})
	state1, _ := stub.GetState(compositeKey1)
	var ad1 models.Admin
	json.Unmarshal([]byte(state1), &ad1)

	// Create a new Proposal - automatically fail if not succeess
	rs4 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateProposal"), []byte("Update Money")})

	// The invokeFunction returns ProposalID key
	var r4 InvokeResponse
	json.Unmarshal([]byte(rs4), &r4)

	// get ProposalID
	compositeKey4, _ := stub.CreateCompositeKey(models.PROPOSALTABLE, []string{r4.Rows})
	state4, _ := stub.GetState(compositeKey4)
	var pr models.Proposal
	json.Unmarshal([]byte(state4), &pr)

	//call CreateQuorum with signature fail
	rs5 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signatureFail, []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var r5 InvokeResponse
	json.Unmarshal([]byte(rs5), &r5)

	fmt.Println("Invoke Quorum signature fail: ", r5.Msg)
	// check err return Fail verify Signature
	assert.Equal(t, "AKC0008", r5.Status)

	//call CreateQuorum with ProposalID not exist
	rs6 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature1, []byte(ad1.AdminID), []byte(ProposalIDFail)})
	var r6 InvokeResponse
	json.Unmarshal([]byte(rs6), &r6)

	fmt.Println("Invoke Quorum return ProposalID not exist: ", r6.Msg)
	//check err return Fail verify Proposal ID not exis
	assert.Equal(t, "AKC0013", r6.Status)

	//call CreateQuorum with AdminID not exist
	rs7 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature1, []byte(AdminIDFail), []byte(pr.ProposalID)})
	var r7 InvokeResponse
	json.Unmarshal([]byte(rs7), &r7)

	fmt.Println("Invoke Quorum return AdminID not exist: ", r7.Msg)
	//check err return Fail verify Admin ID not exist
	assert.Equal(t, "AKC0014", r7.Status)

	//test quorum Only signed once
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature1, []byte(ad1.AdminID), []byte(pr.ProposalID)})

	rs8 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature1, []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var r8 InvokeResponse
	json.Unmarshal([]byte(rs8), &r8)

	compositeKey8, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r8.Rows})
	state8, _ := stub.GetState(compositeKey8)
	var qr8 models.Quorum
	json.Unmarshal([]byte(state8), &qr8)

	//check err return Fail verify Only signed once
	assert.Equal(t, "AKC0009", r8.Status)
}

func TestCommit(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
	pubKey1, _ := ioutil.ReadFile("./sample/pk1.pem")
	pubKey2, _ := ioutil.ReadFile("./sample/pk2.pem")
	pubKey3, _ := ioutil.ReadFile("./sample/pk3.pem")
	signature1, _ := ioutil.ReadFile("./sample/signature1.txt")
	signature2, _ := ioutil.ReadFile("./sample/signature2.txt")
	signature3, _ := ioutil.ReadFile("./sample/signature3.txt")

	// Create a new Admin - automatically fail if not succeess
	rs1 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin1"), []byte(pubKey1)})
	// The invokeFunction returns adminID key
	var r1 InvokeResponse
	json.Unmarshal([]byte(rs1), &r1)
	// Check if the created admin exist in the ledger
	compositeKey1, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r1.Rows})
	state1, _ := stub.GetState(compositeKey1)
	var ad1 models.Admin
	json.Unmarshal([]byte(state1), &ad1)

	rs2 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin2"), []byte(pubKey2)})
	// The invokeFunction returns adminID key
	var r2 InvokeResponse
	json.Unmarshal([]byte(rs2), &r2)
	// Check if the created admin exist in the ledger
	compositeKey2, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r2.Rows})
	state2, _ := stub.GetState(compositeKey2)
	var ad2 models.Admin
	json.Unmarshal([]byte(state2), &ad2)

	rs3 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin3"), []byte(pubKey3)})
	// The invokeFunction returns adminID key
	var r3 InvokeResponse
	json.Unmarshal([]byte(rs3), &r3)
	// Check if the created admin exist in the ledger
	compositeKey3, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r3.Rows})
	state3, _ := stub.GetState(compositeKey3)
	var ad3 models.Admin
	json.Unmarshal([]byte(state3), &ad3)

	// Create a new Proposal - automatically fail if not succeess
	rs4 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateProposal"), []byte("Update Money")})

	// The invokeFunction returns ProposalID key
	var r4 InvokeResponse
	json.Unmarshal([]byte(rs4), &r4)

	// Check if the created Proposal exist in the ledger
	compositeKey4, _ := stub.CreateCompositeKey(models.PROPOSALTABLE, []string{r4.Rows})
	state4, _ := stub.GetState(compositeKey4)
	var pr models.Proposal
	json.Unmarshal([]byte(state4), &pr)

	rs5 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature1, []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var r5 InvokeResponse
	json.Unmarshal([]byte(rs5), &r5)
	compositeKey5, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r5.Rows})
	state5, _ := stub.GetState(compositeKey5)
	var qr1 models.Quorum
	json.Unmarshal([]byte(state5), &qr1)

	rs6 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature2, []byte(ad2.AdminID), []byte(pr.ProposalID)})
	var r6 InvokeResponse
	json.Unmarshal([]byte(rs6), &r6)
	compositeKey6, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r6.Rows})
	state6, _ := stub.GetState(compositeKey6)
	var qr2 models.Quorum
	json.Unmarshal([]byte(state6), &qr2)

	rs7 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature3, []byte(ad3.AdminID), []byte(pr.ProposalID)})
	var r7 InvokeResponse
	json.Unmarshal([]byte(rs7), &r7)

	compositeKey7, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r7.Rows})
	state7, _ := stub.GetState(compositeKey7)
	var qr3 models.Quorum
	json.Unmarshal([]byte(state7), &qr3)

	commitRs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var commitRp InvokeResponse
	json.Unmarshal([]byte(commitRs), &commitRp)

	compositeKey8, _ := stub.CreateCompositeKey(models.COMMITTABLE, []string{commitRp.Rows})
	state8, _ := stub.GetState(compositeKey8)
	var commit models.Commit
	json.Unmarshal([]byte(state8), &commit)

	//confirm create commit success
	assert.Equal(t, "200", commitRp.Status)
	//confirm data in commit struct
	assert.Equal(t, ad1.AdminID, commit.AdminID)
	assert.Equal(t, pr.ProposalID, commit.ProposalID)
	assert.Equal(t, "Verify", commit.Status)
	assert.Contains(t, commit.QuorumList, qr1.QuorumID)
	assert.Contains(t, commit.QuorumList, qr2.QuorumID)
	assert.Contains(t, commit.QuorumList, qr3.QuorumID)

}

func TestCommit_Err(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
	pubKey1, _ := ioutil.ReadFile("./sample/pk1.pem")
	pubKey2, _ := ioutil.ReadFile("./sample/pk2.pem")
	pubKey3, _ := ioutil.ReadFile("./sample/pk3.pem")
	signature1, _ := ioutil.ReadFile("./sample/signature1.txt")
	signature2, _ := ioutil.ReadFile("./sample/signature2.txt")
	signature3, _ := ioutil.ReadFile("./sample/signature3.txt")

	// Create a new Admin - automatically fail if not succeess
	rs1 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin1"), []byte(pubKey1)})
	// The invokeFunction returns adminID key
	var r1 InvokeResponse
	json.Unmarshal([]byte(rs1), &r1)
	// Check if the created admin exist in the ledger
	compositeKey1, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r1.Rows})
	state1, _ := stub.GetState(compositeKey1)
	var ad1 models.Admin
	json.Unmarshal([]byte(state1), &ad1)

	rs2 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin2"), []byte(pubKey2)})
	// The invokeFunction returns adminID key
	var r2 InvokeResponse
	json.Unmarshal([]byte(rs2), &r2)
	// Check if the created admin exist in the ledger
	compositeKey2, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r2.Rows})
	state2, _ := stub.GetState(compositeKey2)
	var ad2 models.Admin
	json.Unmarshal([]byte(state2), &ad2)

	rs3 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin3"), []byte(pubKey3)})
	// The invokeFunction returns adminID key
	var r3 InvokeResponse
	json.Unmarshal([]byte(rs3), &r3)
	// Check if the created admin exist in the ledger
	compositeKey3, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r3.Rows})
	state3, _ := stub.GetState(compositeKey3)
	var ad3 models.Admin
	json.Unmarshal([]byte(state3), &ad3)

	// Create a new Proposal - automatically fail if not succeess
	rs4 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateProposal"), []byte("Update Money")})

	// The invokeFunction returns ProposalID key
	var r4 InvokeResponse
	json.Unmarshal([]byte(rs4), &r4)

	// Check if the created Proposal exist in the ledger
	compositeKey4, _ := stub.CreateCompositeKey(models.PROPOSALTABLE, []string{r4.Rows})
	state4, _ := stub.GetState(compositeKey4)
	var pr models.Proposal
	json.Unmarshal([]byte(state4), &pr)

	commitRs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var commitRp InvokeResponse
	json.Unmarshal([]byte(commitRs), &commitRp)
	//Proposal Commit not exist
	assert.Equal(t, "AKC0012", commitRp.Status)

	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature1, []byte(ad1.AdminID), []byte(pr.ProposalID)})
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature2, []byte(ad2.AdminID), []byte(pr.ProposalID)})

	commitRs1 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var commitRp1 InvokeResponse
	json.Unmarshal([]byte(commitRs1), &commitRp1)
	// 2 quorum Not Enough Quorum
	assert.Equal(t, "AKC0010", commitRp1.Status)

	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature3, []byte(ad3.AdminID), []byte(pr.ProposalID)})

	//commit OK
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad1.AdminID), []byte(pr.ProposalID)})
	//commit duplicate with same admin -> fail
	commitRs2 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var commitRp2 InvokeResponse
	json.Unmarshal([]byte(commitRs2), &commitRp2)
	//Only Commit once!
	assert.Equal(t, "AKC0011", commitRp2.Status)

	//commit duplicate with other admin -> fail
	commitRs3 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad2.AdminID), []byte(pr.ProposalID)})
	var commitRp3 InvokeResponse
	json.Unmarshal([]byte(commitRs3), &commitRp3)
	//Only Commit once!
	assert.Equal(t, "AKC0011", commitRp3.Status)
}
