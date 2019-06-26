package example

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	// Setup mockextend
	cc := new(Chaincode)
	stub := util.NewMockStubExtend(shim.NewMockStub("examp", cc), cc)

	walletAddress := "user0"
	tokenAmount := "100"
	// Create a new User - automatically fail if not succeess
	fmt.Println("Invoke Create User ")
	rs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("createUser"), []byte(walletAddress), []byte(tokenAmount)})

	// The invokeFunction returns wallet address
	var r InvokeResponse
	json.Unmarshal([]byte(rs), &r)

	// Check if the created user exist in the ledger
	compositeKey, _ := stub.CreateCompositeKey(USERTABLE, []string{r.Rows})
	state, _ := stub.GetState(compositeKey)
	var usr User
	json.Unmarshal([]byte(state), &usr)

	// Check if the created user information is correct
	fmt.Println("WalletAddress: ", usr.WalletAddress)
	assert.Equal(t, walletAddress, usr.WalletAddress)
	assert.Equal(t, tokenAmount, fmt.Sprintf("%v", usr.Amount))

	// After pass create user, create second user.
	walletAddress1 := "user1"
	tokenAmount1 := "100"
	fmt.Println("Invoke Create User 2")
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("createUser"), []byte(walletAddress1), []byte(tokenAmount1)})

	// Test balance transfer
	tokenTransfer := "10"
	fmt.Println("Invoke Transfer Token")
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("transferToken"), []byte(walletAddress), []byte(walletAddress1), []byte(tokenTransfer)})

	// Update balance transfer
	fmt.Println("Invoke Update User Balance")
	util.MockInvokeTransaction(t, stub, [][]byte{[]byte("updateUserBalance"), []byte("PRUNE_SAFE")})

	// Get user balance
	fmt.Println("Query get User Balance")
	user1Token := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("getUserToken"), []byte(walletAddress)})
	user2Token := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("getUserToken"), []byte(walletAddress1)})

	var r1 InvokeResponse
	var r2 InvokeResponse
	json.Unmarshal([]byte(user1Token), &r1)
	json.Unmarshal([]byte(user2Token), &r2)

	assert.Equal(t, "90", r1.Rows)
	assert.Equal(t, "110", r2.Rows)
}
