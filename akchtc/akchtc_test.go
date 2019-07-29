package akchtc

import (
	"fmt"
	"strconv"
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

func checkInvokeFail(t *testing.T, stub *shim.MockStub, args [][]byte) (bool, string) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		return true, string(res.Message)
	}
	return false, string(res.Payload)
}

func TestAkcHighThroughtput(t *testing.T) {
	cc := new(SampleChaincode)
	stub := shim.NewMockStub("akchihi", cc)

	// Test Case success
	checkInvoke(t, stub, [][]byte{[]byte("insert"), []byte("Merchant"), []byte("1234567890"), []byte("100"), []byte("OP_ADD")})
	checkInvoke(t, stub, [][]byte{[]byte("insert"), []byte("Merchant"), []byte("0987654321"), []byte("50"), []byte("OP_SUB")})
	checkInvoke(t, stub, [][]byte{[]byte("insert"), []byte("Merchant"), []byte("0987654321"), []byte("25"), []byte("OP_ADD")})
	checkInvoke(t, stub, [][]byte{[]byte("insert"), []byte("Merchant"), []byte("1234567890"), []byte("99"), []byte("OP_SUB")})
	checkInvoke(t, stub, [][]byte{[]byte("insert"), []byte("Merchant"), []byte("88662233"), []byte("50"), []byte("OP_ADD")})

	res2 := checkInvoke(t, stub, [][]byte{[]byte("get"), []byte("Merchant"), []byte("1234567890")})
	res21 := checkInvoke(t, stub, [][]byte{[]byte("get"), []byte("Merchant")})

	res3 := checkInvoke(t, stub, [][]byte{[]byte("prune"), []byte("Merchant"), []byte("1234567890"), []byte("PRUNE_SAFE")})
	// // Prune all with namespace
	res31 := checkInvoke(t, stub, [][]byte{[]byte("prune"), []byte("Merchant"), []byte("PRUNE_SAFE")})
	res32 := checkInvoke(t, stub, [][]byte{[]byte("prune"), []byte("Merchant"), []byte("PRUNE_SAFE")})

	res4 := checkInvoke(t, stub, [][]byte{[]byte("prune"), []byte("Merchant"), []byte("0987654321"), []byte("PRUNE_FAST")})

	res5 := checkInvoke(t, stub, [][]byte{[]byte("delete"), []byte("Merchant"), []byte("1234567890")})

	// // Check output after add 100 to HTC
	s, err1 := strconv.ParseFloat(res2, 64)
	final, err2 := strconv.ParseFloat("1", 64)
	if err1 == nil && err2 == nil && s != final {
		t.Errorf("Inaccurate data, assert: %f, but response: %f", final, s)
	}

	// //============== Test case Insert failure
	// case 1: missing args
	insertFail1, insertMsg1 := checkInvokeFail(t, stub, [][]byte{[]byte("insert"), []byte("User"), []byte("1"), []byte("100")})
	if insertMsg1 != "Incorrect number of arguments, expecting 4" {
		t.Errorf("Check insert fail: Assert expect not nil, but response %v.", insertFail1)
	}

	// // case 2: operation unrecognized
	insertFail2, insertMsg2 := checkInvokeFail(t, stub, [][]byte{[]byte("insert"), []byte("User"), []byte("2"), []byte("100"), []byte("?")})
	if insertMsg2 != "Operator ? is unrecognized" {
		t.Errorf("Check insert fail: Assert expect pperator is unrecognized, but response %v", insertFail2)
	}

	// // case 3: value not a number
	insertFail3, insertMsg3 := checkInvokeFail(t, stub, [][]byte{[]byte("insert"), []byte("User"), []byte("2"), []byte("abc"), []byte("OP_ADD")})

	if !insertFail3 {
		t.Errorf("Insert with value `abc` is fail, but response %v", insertFail3)
	}

	// //============== Test case Get failure
	// case 1: missing args
	getFail1, getMsg1 := checkInvokeFail(t, stub, [][]byte{[]byte("get")})
	if getMsg1 != "Incorrect number of arguments, expecting 1" {
		t.Errorf("Check get fail: Assert expect 1 args, but response %v", getFail1)
	}

	// case 2: No variable by name exists
	getFail2, getMsg2 := checkInvokeFail(t, stub, [][]byte{[]byte("get"), []byte("ahihi"), []byte("123")})
	if getMsg2 != "No variable by the name ahihi exists" {
		t.Errorf("Check get fail: Assert expect name `ahihi` not exists, but response %v", getFail2)
	}

	// //============== Test case Prune failure
	// case 1: missing args
	pruneFail1, pruneMsg1 := checkInvokeFail(t, stub, [][]byte{[]byte("prune"), []byte("Merchant")})
	if pruneMsg1 != "Incorrect number of arguments, expecting 2" {
		t.Errorf("Check prune fail: Assert expect 2 args, but response success. %v", pruneFail1)
	}

	// case 2: No variable by name exists
	pruneFail2, pruneMsg2 := checkInvokeFail(t, stub, [][]byte{[]byte("prune"), []byte("ahihi"), []byte("123"), []byte("PRUNE_SAFE")})
	if !pruneFail2 {
		t.Errorf("Check prune fail: Assert expect name `ahihi` not exists and return prune FALSE, but response %v", pruneFail2)
	}

	// case 2: Prune type not supported
	pruneFail3, pruneMsg3 := checkInvokeFail(t, stub, [][]byte{[]byte("prune"), []byte("ahihi"), []byte("123"), []byte("PRUNE_SLOW")})
	if pruneMsg3 != "Prune type PRUNE_SLOW is not supported" {
		t.Errorf("Check prune fail: Assert expect type `PRUNE_SLOW` not supported, but response %v", pruneFail3)
	}

	// //############# Test response #################
	fmt.Printf("Get success response: %v\n", res2)
	fmt.Printf("Get success response: %v\n", res21)
	fmt.Printf("Prune success response: %v\n", res3)
	fmt.Printf("Prune success response: %v\n", res31)
	fmt.Printf("Prune success response: %v\n", res32)
	fmt.Printf("Prune success response: %v\n", res4)
	fmt.Printf("Delete success response: %v\n", res5)

	fmt.Printf("Insert case 1 fail response: %v\n", insertMsg1)
	fmt.Printf("Insert case 2 fail response: %v\n", insertMsg2)
	fmt.Printf("Insert case 3 fail response: %v\n", insertMsg3)

	fmt.Printf("Get case 1 fail response: %v\n", getMsg1)
	fmt.Printf("Get case 2 fail response: %v\n", getMsg2)

	fmt.Printf("Prune case 1 fail response: %v\n", pruneMsg1)
	fmt.Printf("Prune case 2 fail response: %v\n", pruneMsg2)
	fmt.Printf("Prune case 3 fail response: %v\n", pruneMsg3)
}
