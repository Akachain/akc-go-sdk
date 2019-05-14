package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkByID111(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) string {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		return string(res.Payload)
	}
	return string(res.Payload)
}

func checkGetALL(t *testing.T, stub *shim.MockStub, chaincodename string) []byte {
	res := stub.MockInvoke("1", [][]byte{[]byte(chaincodename)})
	if res.Status != shim.OK {
		fmt.Println(chaincodename, string(res.Payload))
		t.FailNow()
	}
	return res.Payload
}

func checkByID(t *testing.T, stub *shim.MockStub, args [][]byte) []byte {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		t.FailNow()
		return res.Payload
	}
	return res.Payload
}

func TestMainchain(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("mainchain", cc)
	AdminID := "ntienbo"

	// Check get CreateProposal
	rs := checkInvoke(t, stub, [][]byte{[]byte("CreateProposal"), []byte("1")})

	rs5 := checkInvoke(t, stub, [][]byte{[]byte("CreateAdmin"), []byte("Admin"), []byte("pulic key")})

	rs3 := checkByID(t, stub, [][]byte{[]byte("GetAdminByID"), []byte("Admin")})

	rs4 := checkGetALL(t, stub, "GetAllProposal")

	rs6 := checkGetALL(t, stub, "GetAllAdmin")

	ProposalID := "d"
	sig := "MEYCIQC7vKLzjw43HJ/9SqxUzZtfdBIdFks7qiIXiHitu8uqqQIhAKXNwpBDuWquPE/00l8isa6rh85ZYYf+dgb1khSqNr7O"
	rs2 := checkInvoke(t, stub, [][]byte{[]byte("CreateQuorum"), []byte(sig), []byte(AdminID), []byte(ProposalID)})

	fmt.Printf("rs: %v", rs)
	fmt.Printf("rs2: %v", rs2)
	fmt.Printf("rs3: %v", rs3)
	fmt.Printf("rs4: %v", rs4)
	fmt.Printf("rs5: %v", rs5)
	fmt.Printf("rs6: %v", rs6)

	// check checkByID
}
