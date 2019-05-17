package controllers

import (
	"fmt"

	"github.com/rs/xid"

	"github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	. "github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"
)

type Admin models.Admin

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

	err := Createdata(stub, models.ADMINTABLE, []string{AdminID}, &Admin{AdminID: AdminID, Name: Name, PublicKey: Publickey})

	if err != nil {
		resErr := common.ResponseError{common.ERR5, fmt.Sprintf("%s %s", common.ResCodeDict[common.ERR5], "")}
		return common.RespondError(resErr)
	}
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], AdminID}
	return common.RespondSuccess(resSuc)
}

//UpdateAdmin
func (admin *Admin) UpdateAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	// get admin information
	var admin_tmp Admin
	AdminID := args[0]
	admin_rs, err := Getdatabyid(stub, AdminID, models.ADMINTABLE)
	if err != nil {
		resErr := common.ResponseError{common.ERR4, fmt.Sprintf("%s %s", common.ResCodeDict[common.ERR4], err.Error())}
		return common.RespondError(resErr)
	}

	//Find Field Need update
	mapstructure.Decode(admin_rs, &admin_tmp)
	if args[1] == "Name" {
		admin_tmp.Name = args[2]
	} else if args[1] == "PublicKey" {
		admin_tmp.PublicKey = args[2]
	}

	err = Changeinfo(stub, models.ADMINTABLE, []string{AdminID}, &Admin{AdminID: AdminID, Name: admin_tmp.Name, PublicKey: admin_tmp.PublicKey})
	if err != nil {
		//Overwrite fail
		resErr := common.ResponseError{common.ERR5, fmt.Sprintf("%s %s", common.ResCodeDict[common.ERR5], err.Error())}
		return common.RespondError(resErr)
	}
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], "[]"}
	return common.RespondSuccess(resSuc)
}

//GetAdminByID
func (admin *Admin) GetAdminByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	DataID := args[0]
	res := GetDataByID(stub, DataID, admin, models.ADMINTABLE)
	return res
}

//GetAllAdmin
func (admin *Admin) GetAllAdmin(stub shim.ChaincodeStubInterface) pb.Response {
	res := GetAllData(stub, admin, models.ADMINTABLE)
	return res
}

// ------------------- //
