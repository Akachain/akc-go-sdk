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

func TestDBHandle(t *testing.T) {
	a := util.NewCouchDBHandler()
	Logger.Debug(a)
}

func TestAdmin(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
	admin := "Admin1"
	pubKey, _ := ioutil.ReadFile("./sample/pk.pem")
	pk := base64.StdEncoding.EncodeToString(pubKey)

	// Create a new Admin - automatically fail if not succeess
	fmt.Println("Invoke CreateAdmin ", admin)
	rs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("CreateAdmin"), []byte(admin), []byte(pk)})

	// The invokeFunction returns adminID key
	var r InvokeResponse
	json.Unmarshal([]byte(rs), &r)

	// Check if the created admin exist in the ledger
	compositeKey, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r.Rows})
	state, _ := stub.GetState(compositeKey)
	var ad models.Admin
	json.Unmarshal([]byte(state), &ad)

	// Check if the created admin information is correct
	assert.Equal(t, admin, ad.Name)
	assert.Equal(t, pk, ad.PublicKey)
}

func TestQuorum(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("hstx", cc), cc)
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
