package controllers

import (
	"fmt"

	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/models"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"
)

type Admin models.Admin

//CreateAdmin adds an admin document that contain the AdminID and his Public Key
func (admin *Admin) CreateAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 2)
	AdminID := fmt.Sprintf("%s_%s", stub.GetTxID(), args[0])
	Name := args[0]
	Publickey := args[1]
	Status := "Active"

	err := util.Createdata(stub, models.ADMINTABLE, []string{AdminID}, &Admin{AdminID: AdminID, Name: Name, PublicKey: Publickey, Status: Status})
	if err != nil {
		resErr := ResponseError{ERR5, fmt.Sprintf("%s %s", ResCodeDict[ERR5], "")}
		return RespondError(resErr)
	}
	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], AdminID}
	return RespondSuccess(resSuc)
}

//UpdateAdmin Edit Status to "Delete" if you want to delete admin
func (admin *Admin) UpdateAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 3)
	//get admin information
	var admin_tmp Admin
	AdminID := args[0]
	admin_rs, err := util.Getdatabyid(stub, AdminID, models.ADMINTABLE)
	if err != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s", ResCodeDict[ERR4], err.Error())}
		return RespondError(resErr)
	}

	//Find Field Need update
	mapstructure.Decode(admin_rs, &admin_tmp)
	if args[1] == "Name" {
		admin_tmp.Name = args[2]
	} else if args[1] == "PublicKey" {
		admin_tmp.PublicKey = args[2]
	} else if args[1] == "Status" {
		admin_tmp.Status = args[2]
	}

	err = util.Changeinfo(stub, models.ADMINTABLE, []string{AdminID}, &Admin{AdminID: AdminID, Name: admin_tmp.Name, PublicKey: admin_tmp.PublicKey, Status: admin_tmp.Status})
	if err != nil {
		//Overwrite fail
		resErr := ResponseError{ERR5, fmt.Sprintf("%s %s", ResCodeDict[ERR5], err.Error())}
		return RespondError(resErr)
	}
	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], "[]"}
	return RespondSuccess(resSuc)
}

//GetAdminByID
func (admin *Admin) GetAdminByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 1)

	DataID := args[0]
	res := util.GetDataByID(stub, DataID, admin, models.ADMINTABLE)
	return res
}

//GetAllAdmin
func (admin *Admin) GetAllAdmin(stub shim.ChaincodeStubInterface) pb.Response {
	res := util.GetAllData(stub, admin, models.ADMINTABLE)
	return res
}

// ------------------- //
