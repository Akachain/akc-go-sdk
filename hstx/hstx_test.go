package main

import (
	"encoding/base64"
	"encoding/hex"
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
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
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
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
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

	fmt.Println("Invoke Quorum ID 1", qr1.QuorumID)
	fmt.Println("Invoke Quorum 1", r1)

	// Check if the created admin information is correct
	assert.Equal(t, pr_in, pr.Data)
}

func TestQuorum_Err(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
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
	//check err return Fail verify Proposal ID not exis
	assert.Equal(t, "AKC0013", r4.Status)

	//call CreateQuorum with AdminID not exist
	rs5 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(AdminIDFail), []byte(pr.ProposalID)})
	var r5 InvokeResponse
	json.Unmarshal([]byte(rs5), &r5)

	fmt.Println("Invoke Quorum return AdminID not exist: ", r5.Msg)
	//check err return Fail verify Admin ID not exist
	assert.Equal(t, "AKC0014", r5.Status)

	//test quorum Only signed once
	_ = util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(pr.ProposalID)})
	rs7 := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(pr.ProposalID)})
	var r7 InvokeResponse
	json.Unmarshal([]byte(rs7), &r7)

	compositeKey7, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{r7.Rows})
	state7, _ := stub.GetState(compositeKey7)
	var qr7 models.Quorum
	json.Unmarshal([]byte(state7), &qr7)

	//check err return Fail verify Only signed once
	assert.Equal(t, "AKC0009", r7.Status)

}

func TestCommit(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)
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

	_ = util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(pr.ProposalID)})

	_ = util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(pr.ProposalID)})

	_ = util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ad.AdminID), []byte(pr.ProposalID)})

	commitRs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateCommit"), []byte(pr.ProposalID)})
	var commitRp InvokeResponse
	json.Unmarshal([]byte(commitRs), &commitRp)

	compositeKey6, _ := stub.CreateCompositeKey(models.QUORUMTABLE, []string{commitRp.Rows})
	state6, _ := stub.GetState(compositeKey6)
	var commit models.Commit
	json.Unmarshal([]byte(state6), &commit)

	fmt.Printf("Invoke Commit ID: %v \n", commit.CommitID)
	fmt.Printf("Invoke Commit: %v \n", commitRp)

}

func TestQuorumLong(t *testing.T) {
	// Setup mockextend - split this to util - TODO
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	db := util.NewCouchDBHandler("hstx-test")
	stub.SetCouchDBConfiguration(db)

	// Sample data
	admin := "Admin1"
	pubKey, _ := ioutil.ReadFile("./sample/pk.pem")
	pk := base64.StdEncoding.EncodeToString(pubKey)

	// Create a new Admin - automatically fail if not succeess
	fmt.Println("Invoke CreateAdmin ", admin)
	res := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte(admin), []byte(pk)})
	var ra InvokeResponse
	json.Unmarshal([]byte(res), &ra)

	// Create a new Proposal - automatically fail if not succeess
	fmt.Println("Invoke CreateProposal ", admin)
	res = util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateProposal"), []byte("Secure Transaction")})
	var rp InvokeResponse
	json.Unmarshal([]byte(res), &rp)

	// Create a new Quorum - automatically fail if not succeess
	fmt.Println("Invoke CreateQuorum ", admin)
	signature, _ := hex.DecodeString("250cefebe48f40aa50c369d5842f8bab79223467226cc4cdc573a09afbcf668317d24cd8e10a5e1a7b3b701126c02341714131e2477425fac95b10df987141c241f3cf4db04c8f536abbfb01f67db056e27994c55545d77f8293505bb35437b23ea4d178b77b6f6aa9994292b2d8eb3947d5f9a79e1730d96152612650c8072ffb639a8d92c4dda146d8fa248fd559199829c8d6eb7bd5449a3f162e338daf2ff199671b460a81ea42d1146fbeb0514cdc42c07717723fff397fda34c93e12a399947eea8bc1da98d872a8d9f8c87ef541970aa0d7774318880cd17c4578781bbb65c670e76ebd675ca79f653449840e454d3847581c5c865f79a6678ddc2ae7")
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateQuorum"), signature, []byte(ra.Rows), []byte(rp.Rows)})
}
