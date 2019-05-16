package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkCallFuncInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) string {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		return string(res.Payload)
	}
	return string(res.Payload)
}
func checkCallFuncQuerry(t *testing.T, stub *shim.MockStub, args [][]byte) string {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		t.FailNow()
		return string(res.Payload)
	}
	t.FailNow()
	return string(res.Payload)
}

func TestMainchain(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("mainchain", cc)

	//	rs5 := checkCallFuncInvoke(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin"), []byte("pulic key")})
	rs7 := checkCallFuncInvoke(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin1"), []byte("pulic key")})
	rs7 = checkCallFuncInvoke(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin21"), []byte("pulic key")})

	//rs3 := checkByID(t, stub, [][]byte{[]byte("GetAdminByID1"), []byte("aa"), []byte(models.ADMINTABLE)})

	//rs4 := checkGetALL(t, stub, "GetAllProposal")

	//	rs6 = checkCallFuncQuerry(t, stub, [][]byte{[]byte("GetAllAdmin"), []byte("")})
	rs6 := checkCallFuncQuerry(t, stub, [][]byte{[]byte("GetAdminByID"), []byte("a")})

	//fmt.Printf("rs3: %v", rs3)
	//	fmt.Printf("rs4: %v", rs4)
	//	fmt.Printf("rs5: %v", rs5)
	fmt.Printf("rs6: %v", rs6)
	fmt.Printf("rs6: %v", rs7)
	//	fmt.Printf("rs6: %v", rs8)

	// check checkByID
}
