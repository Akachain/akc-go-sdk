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
	rs := util.MockInitTransaction(t, stub, [][]byte{[]byte("")})
	assert.Equal(t, "", rs)
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

	// The invokeFunction returns adminID key
	var r InvokeResponse
	json.Unmarshal([]byte(rs), &r)

	// Check if the created admin exist in the ledger
	compositeKey, _ := stub.CreateCompositeKey(models.PROPOSALTABLE, []string{r.Rows})
	state, _ := stub.GetState(compositeKey)
	var pr models.Proposal
	json.Unmarshal([]byte(state), &pr)

	// Check if the created admin information is correct
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

	// Create a new Proposal - automatically fail if not succeess
	rs2 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateProposal"), []byte("Update Money")})

	// The invokeFunction returns adminID key
	var r1 InvokeResponse
	json.Unmarshal([]byte(rs1), &r1)

	// Check if the created admin exist in the ledger
	compositeKey1, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r1.Rows})
	state1, _ := stub.GetState(compositeKey1)
	var ad models.Admin
	json.Unmarshal([]byte(state1), &ad)

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

	// Check if the created admin information is correct
	assert.Equal(t, "200", r3.Status)
}

func TestQuorum_Err(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
	pubKey1, _ := ioutil.ReadFile("./sample/pk1.pem")
	pubKey2, _ := ioutil.ReadFile("./sample/pk2.pem")
	pubKey3, _ := ioutil.ReadFile("./sample/pk3.pem")
	signature1, _ := ioutil.ReadFile("./sample/signature1.txt")
	// signature2, _ := ioutil.ReadFile("./sample/signature2.txt")
	// signature3, _ := ioutil.ReadFile("./sample/signature3.txt")
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

	// Create a new Admin - automatically fail if not succeess
	rs2 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin2"), []byte(pubKey2)})
	// The invokeFunction returns adminID key
	var r2 InvokeResponse
	json.Unmarshal([]byte(rs2), &r2)

	// Check if the created admin exist in the ledger
	compositeKey2, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r2.Rows})
	state2, _ := stub.GetState(compositeKey2)
	var ad2 models.Admin
	json.Unmarshal([]byte(state2), &ad2)

	// Create a new Admin - automatically fail if not succeess
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
	// util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature2, []byte(ad2.AdminID), []byte(pr.ProposalID)})
	// util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature3, []byte(ad3.AdminID), []byte(pr.ProposalID)})

	rs8 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature1, []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var r8 InvokeResponse
	json.Unmarshal([]byte(rs8), &r8)

	compositeKey8, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r8.Rows})
	state8, _ := stub.GetState(compositeKey8)
	var qr8 models.Quorum
	json.Unmarshal([]byte(state8), &qr8)

	//check err return Fail verify Only signed once
	assert.Equal(t, "AKC00091", r8.Status)
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

	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature1, []byte(ad1.AdminID), []byte(pr.ProposalID)})
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature2, []byte(ad2.AdminID), []byte(pr.ProposalID)})
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature3, []byte(ad3.AdminID), []byte(pr.ProposalID)})

	commitRs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var commitRp InvokeResponse
	json.Unmarshal([]byte(commitRs), &commitRp)

	compositeKey6, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{commitRp.Rows})
	state6, _ := stub.GetState(compositeKey6)
	var commit models.Commit
	json.Unmarshal([]byte(state6), &commit)

	fmt.Printf("Invoke Commit ID: %v \n", commit.CommitID)
	fmt.Printf("Invoke Commit: %v \n", commitRp)
	assert.Equal(t, "200", commitRp.Status)
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

	commitRs2 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var commitRp2 InvokeResponse
	json.Unmarshal([]byte(commitRs2), &commitRp2)
	//Proposal Commit not exist
	assert.Equal(t, "AKC0012", commitRp2.Status)

	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature1, []byte(ad1.AdminID), []byte(pr.ProposalID)})
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature2, []byte(ad2.AdminID), []byte(pr.ProposalID)})

	commitRs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var commitRp InvokeResponse
	json.Unmarshal([]byte(commitRs), &commitRp)
	//Not Enough Quorum
	assert.Equal(t, "AKC0010", commitRp.Status)

	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature3, []byte(ad3.AdminID), []byte(pr.ProposalID)})

	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad1.AdminID), []byte(pr.ProposalID)})
	commitRs1 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(ad1.AdminID), []byte(pr.ProposalID)})
	var commitRp1 InvokeResponse
	json.Unmarshal([]byte(commitRs1), &commitRp1)
	//Only Commit once!
	assert.Equal(t, "AKC0011", commitRp1.Status)

}
