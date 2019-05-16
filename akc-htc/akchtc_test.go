package akchtc

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) string {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		return string(res.Payload)
	}
	return string(res.Payload)
}

func TestAkcHighThroughtput(t *testing.T) {
	cc := new(SampleChaincode)
	stub := shim.NewMockStub("akchihi", cc)

	res1 := checkInvoke(t, stub, [][]byte{[]byte("insert"), []byte("Merchant"), []byte("1234567890"), []byte("100"), []byte("+")})

	res2 := checkInvoke(t, stub, [][]byte{[]byte("get"), []byte("Merchant"), []byte("1234567890")})

	res3 := checkInvoke(t, stub, [][]byte{[]byte("prune"), []byte("Merchant"), []byte("1234567890"), []byte("PRUNE_SAFE")})

	res4 := checkInvoke(t, stub, [][]byte{[]byte("prune"), []byte("Merchant"), []byte("1234567890"), []byte("PRUNE_FAST")})

	res5 := checkInvoke(t, stub, [][]byte{[]byte("delete"), []byte("Merchant"), []byte("1234567890")})

	fmt.Printf("Insert response: %v", res1)
	fmt.Println(".")
	fmt.Printf("Get response: %v", res2)
	fmt.Println(".")
	fmt.Printf("Prune response: %v", res3)
	fmt.Println(".")
	fmt.Printf("Prune response: %v", res4)
	fmt.Println(".")
	fmt.Printf("Delete response: %v", res5)
}
