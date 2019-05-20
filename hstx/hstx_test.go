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
	AdminID := "ntienbo"

	//rs5 := checkCallFuncInvoke(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin"), []byte("pulic key")})
	rs7 := checkCallFuncInvoke(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin1"), []byte("pulic key")})

	rs7 = checkCallFuncInvoke(t, stub, [][]byte{[]byte("UpdateAdmin"), []byte("a"), []byte("Name"), []byte("Newnew")})

	ProposalID := "a"
	sig := "MEYCIQC7vKLzjw43HJ/9SqxUzZtfdBIdFks7qiIXiHitu8uqqQIhAKXNwpBDuWquPE/00l8isa6rh85ZYYf+dgb1khSqNr7O"
	rs7 = checkCallFuncInvoke(t, stub, [][]byte{[]byte("CreateQuorum"), []byte(sig), []byte(AdminID), []byte(ProposalID)})

	rs6 := checkCallFuncQuerry(t, stub, [][]byte{[]byte("GetAllAdmin"), []byte("")})
	rs6 = checkCallFuncQuerry(t, stub, [][]byte{[]byte("GetAdminByID"), []byte("a")})

	//fmt.Printf("rs3: %v", rs3)
	//fmt.Printf("rs4: %v", rs4)
	//fmt.Printf("rs5: %v", rs5)
	fmt.Printf("rs6: %v", rs6)
	fmt.Printf("rs6: %v", rs7)
	//fmt.Printf("rs6: %v", rs8)

	//check checkByID
}
