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

func checkCallFuncInvoke(t *testing.T, stub *util.MockStubExtend, args [][]byte) string {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		return string(res.Payload)
	}
	return string(res.Payload)
}

func checkCallFuncQuerry(t *testing.T, stub *util.MockStubExtend, args [][]byte) string {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		t.FailNow()
		return string(res.Payload)
	}
	t.FailNow()
	return string(res.Payload)
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
	rs := checkCallFuncInvoke(t, stub, [][]byte{[]byte("CreateAdmin"), []byte(admin), []byte(pk)})

	// The invokeFunction returns adminID key
	var r InvokeResponse
	json.Unmarshal([]byte(rs), &r)

	// Check if the created admin exist in the ledger
	composite_key, _ := stub.CreateCompositeKey(models.ADMINTABLE, []string{r.Rows})
	state, _ := stub.GetState(composite_key)
	var ad models.Admin
	json.Unmarshal([]byte(state), &ad)

	assert.Equal(t, admin, ad.Name)
	assert.Equal(t, pk, ad.PublicKey)
}
