package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Akachain/akc-go-sdk/common"
	hdl "github.com/Akachain/akc-go-sdk/hstx/handler"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode struct
type Chaincode struct {
}

var handler = new(hdl.Handler)

// Init method is called when the Chain code" is instantiated by the blockchain network
func (s *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the chain code
// The calling application program has also specified the particular smart contract function to be called, with arguments
func (s *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	common.Logger.Info("########### Hstx Invoke ###########")

	handler.InitHandler()

	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()

	router := map[string]interface{}{
		"CreateSuperAdmin": handler.SuperAdminHanler.CreateSuperAdmin,
		"UpdateSuperAdmin": handler.SuperAdminHanler.UpdateSuperAdmin,

		"CreateAdmin": handler.AdminHanler.CreateAdmin,
		"UpdateAdmin": handler.AdminHanler.UpdateAdmin,

		"CreateProposal": handler.ProposalHanler.CreateProposal,
		"UpdateProposal": handler.ProposalHanler.UpdateProposal,

		"CreateApproval": handler.ApprovalHanler.CreateApproval,
		"UpdateApproval": handler.ApprovalHanler.UpdateApproval,
	}

	routerType := "invoke"
	return s.route(routerType, router, function, stub, args)
}

// Query callback representing the query of a chaincode
func (s *Chaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	common.Logger.Info("########### Hstx Query ###########")

	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()

	router := map[string]interface{}{
		"GetAllSuperAdmin":  handler.SuperAdminHanler.GetAllSuperAdmin,
		"GetSuperAdminByID": handler.SuperAdminHanler.GetSuperAdminByID,

		"GetAllAdmin":  handler.AdminHanler.GetAllAdmin,
		"GetAdminByID": handler.AdminHanler.GetAdminByID,

		"GetAllProposal":  handler.ProposalHanler.GetAllProposal,
		"GetProposalByID": handler.ProposalHanler.GetProposalByID,

		"GetAllApproval":  handler.ApprovalHanler.GetAllApproval,
		"GetApprovalByID": handler.ApprovalHanler.GetApprovalByID,
	}

	routerType := "query"
	return s.route(routerType, router, function, stub, args)
}

// route func to route the funcName to a hanler's function correspondingly
func (s *Chaincode) route(routerType string, router map[string]interface{}, funcName string, params ...interface{}) pb.Response {
	if router[funcName] == nil {
		if strings.Compare(routerType, "invoke") == 0 {
			return s.Query(params[0].(shim.ChaincodeStubInterface))
		}
		return shim.Error(fmt.Sprintf("[Hstx Chaincode] Invoke not find function " + funcName))
	}

	function := reflect.ValueOf(router[funcName])
	if !function.IsValid() {
		return shim.Error(fmt.Sprintf("Map function is invalid."))
	}

	if len(params) != function.Type().NumIn() {
		return shim.Error(fmt.Sprintf("The number of params is not adapted."))
	}

	paramsValue := make([]reflect.Value, len(params))
	for i, param := range params {
		paramsValue[i] = reflect.ValueOf(param)
	}

	funcResult := reflect.MakeFunc(reflect.TypeOf(router[funcName]), func(paramsValue []reflect.Value) []reflect.Value {
		return function.Call(paramsValue)
	})

	f := funcResult.Interface().(func(shim.ChaincodeStubInterface, []string) pb.Response)
	return f(params[0].(shim.ChaincodeStubInterface), params[1].([]string))
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Chain code
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error creating new Chain code: %s", err)
	}
}
