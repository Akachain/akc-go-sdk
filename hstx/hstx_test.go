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

func TestAdmin(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	adminName := "Admin1"
	pubKey, _ := ioutil.ReadFile("./sample/pk.pem")
	pk := base64.StdEncoding.EncodeToString(pubKey)

	// Create a new Admin - automatically fail if not succeess
	fmt.Println("Invoke CreateAdmin ", adminName)
	rs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin2"), []byte(pk)})

	// The invokeFunction returns adminID key
	var r InvokeResponse
	json.Unmarshal([]byte(rs), &r)

	// Check if the created admin exist in the ledger
	compositeKey, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r.Rows})
	state, _ := stub.GetState(compositeKey)
	var ad models.Admin
	json.Unmarshal([]byte(state), &ad)

	// Check if the created admin information is correct
	fmt.Println("AdminID: ", ad.AdminID)
	assert.Equal(t, adminName, ad.Name)
	assert.Equal(t, pk, ad.PublicKey)

}
func TestProposal(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	pr_data := "Update Money"

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
	assert.Equal(t, pr_data, pr.Data)
}

func TestQuorum(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	pr_in := "Update Money1"
	adminName := "Admin1"
	pubKey, _ := ioutil.ReadFile("./sample/pk.pem")
	signature, _ := ioutil.ReadFile("./sample/signature.txt")

	// Create a new Admin - automatically fail if not succeess
	fmt.Println("Invoke CreateAdmin ", adminName)
	rs1 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte(adminName), []byte(pubKey)})

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

	rs3 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(pr.ProposalID)})
	var r3 InvokeResponse
	json.Unmarshal([]byte(rs3), &r3)

	compositeKey3, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r3.Rows})
	state3, _ := stub.GetState(compositeKey3)
	var qr models.Quorum
	json.Unmarshal([]byte(state3), &qr)

	fmt.Println("Invoke Quorum ID ", qr.QuorumID)
	fmt.Println("Invoke Quorum ", r3)

	// Check if the created admin information is correct
	assert.Equal(t, pr_in, pr.Data)
}

func TestQuorum_Err(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	adminName := "Admin1"
	pubKey, _ := ioutil.ReadFile("./sample/pk.pem")
	signature, _ := ioutil.ReadFile("./sample/signature.txt")
	signatureFail, _ := ioutil.ReadFile("./sample/signatureFail.txt")
	ProposalIDFail := "ProposalID no thing"
	AdminIDFail := "AdminID no thing"

	// Create a new Admin - automatically fail if not succeess
	fmt.Println("Invoke CreateAdmin ", adminName)
	rs1 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte(adminName), []byte(pubKey)})

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

	// get ProposalID
	compositeKey2, _ := stub.CreateCompositeKey(models.PROPOSALTABLE, []string{r2.Rows})
	state2, _ := stub.GetState(compositeKey2)
	var pr models.Proposal
	json.Unmarshal([]byte(state2), &pr)

	//call CreateQuorum with signature fail
	rs3 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signatureFail, []byte(ad.AdminID), []byte(pr.ProposalID)})
	var r3 InvokeResponse
	json.Unmarshal([]byte(rs3), &r3)

	fmt.Println("Invoke Quorum signature fail: ", r3.Msg)
	// check err return Fail verify Signature
	assert.Equal(t, "AKC0008", r3.Status)

	//call CreateQuorum with ProposalID not exist
	rs4 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(ProposalIDFail)})
	var r4 InvokeResponse
	json.Unmarshal([]byte(rs4), &r4)

	fmt.Println("Invoke Quorum return ProposalID not exist: ", r4.Msg)
	//check err return Fail verify
	assert.Equal(t, "AKC0013", r4.Status)

	//call CreateQuorum with AdminID not exist
	rs5 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(AdminIDFail), []byte(pr.ProposalID)})
	var r5 InvokeResponse
	json.Unmarshal([]byte(rs5), &r5)

	fmt.Println("Invoke Quorum return AdminID not exist: ", r5.Msg)
	//check err return Fail verify
	assert.Equal(t, "AKC0014", r5.Status)

	//call CreateQuorum with 4 args
	rs6 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(pr.ProposalID), []byte(pr.ProposalID)})
	var r6 InvokeResponse
	json.Unmarshal([]byte(rs6), &r6)

	fmt.Println("Invoke Quorum return args err: ", r6.Msg)
	//check err return Fail verify
	assert.Equal(t, "AKC00114", r6.Status)
}

func TestCommit(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	pr_in := "Update Money1"
	adminName := "Admin1"
	pubKey, _ := ioutil.ReadFile("./sample/pk.pem")
	signature, _ := ioutil.ReadFile("./sample/signature.txt")

	// Create a new Admin - automatically fail if not succeess
	fmt.Println("Invoke CreateAdmin ", adminName)
	rs1 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte(adminName), []byte(pubKey)})

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

	rs3 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(pr.ProposalID)})
	var r3 InvokeResponse
	json.Unmarshal([]byte(rs3), &r3)

	compositeKey3, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r3.Rows})
	state3, _ := stub.GetState(compositeKey3)
	var qr1 models.Quorum
	json.Unmarshal([]byte(state3), &qr1)

	fmt.Printf("Invoke Quorum1 ID: %v \n", qr1.QuorumID)
	fmt.Printf("Invoke Quorum1: %v \n", r3)

	rs4 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(pr.ProposalID)})
	var r4 InvokeResponse
	json.Unmarshal([]byte(rs4), &r4)

	compositeKey4, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r4.Rows})
	state4, _ := stub.GetState(compositeKey4)
	var qr2 models.Quorum
	json.Unmarshal([]byte(state4), &qr2)

	fmt.Printf("Invoke Quorum2 ID: %v \n", qr2.QuorumID)
	fmt.Printf("Invoke Quorum2: %v \n", r4)

	rs5 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(pr.ProposalID)})
	var r5 InvokeResponse
	json.Unmarshal([]byte(rs5), &r5)

	compositeKey5, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r5.Rows})
	state5, _ := stub.GetState(compositeKey5)
	var qr3 models.Quorum
	json.Unmarshal([]byte(state5), &qr3)

	fmt.Printf("Invoke Quorum3 ID: %v \n", qr3.QuorumID)
	fmt.Printf("Invoke Quorum3: %v \n", r5)

	rs6 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(pr.ProposalID)})
	var r6 InvokeResponse
	json.Unmarshal([]byte(rs6), &r6)

	compositeKey6, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r6.Rows})
	state6, _ := stub.GetState(compositeKey6)
	var commit models.Commit
	json.Unmarshal([]byte(state6), &commit)

	fmt.Printf("Invoke Commit ID: %v \n", commit.CommitID)
	fmt.Printf("Invoke Commit: %v \n", r6)

	// Check if the created admin information is correct
	assert.Equal(t, pr_in, pr.Data)
}
