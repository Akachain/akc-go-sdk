package controllers

import (
	"fmt"

	"github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/rs/xid"
)

//High secure transaction Admin handle
// ------------------- //
//create Admin
func (admin *Admin) CreateAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	AdminID := xid.New().String()
	Name := args[0]
	Publickey := args[1]

	err := create_data_(stub, models.ADMINTABLE, []string{AdminID}, &Admin{AdminID: AdminID, Name: Name, PublicKey: Publickey})

	if err != nil {
		resErr := common.ResponseError{common.ERR5, fmt.Sprintf("%s %s", common.ResCodeDict[common.ERR5], "")}
		return common.RespondError(resErr)
	}
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], AdminID}
	return common.RespondSuccess(resSuc)
}

// ------------------- //
