package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
	"gitlab.com/akachain/akc-go-sdk/common"
	"gitlab.com/akachain/akc-go-sdk/hstx/models"
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

//Get Admin By ID
func (admin *Admin) GetAdminByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	AdminID := args[0]

	rs, err := get_data_byid_(stub, AdminID, models.ADMINTABLE)

	mapstructure.Decode(rs, admin)
	fmt.Printf("Amdin: %v\n", admin)

	bytes, err := json.Marshal(admin)
	if err != nil {
		//Convert Json Fail
		resErr := common.ResponseError{common.ERR3, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	fmt.Printf("Response: %s\n", string(bytes))

	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(bytes)}
	return common.RespondSuccess(resSuc)
}

//Get all Admin
func (admin *Admin) GetAllAdmin(stub shim.ChaincodeStubInterface) pb.Response {
	adminbytes, err := get_all_data_(stub, models.ADMINTABLE)

	admin = new(Admin)
	Adminlist := []*Admin{}

	for row_json_bytes := range adminbytes {
		admin = new(Admin)
		err = json.Unmarshal(row_json_bytes, admin)
		if err != nil {

			resErr := common.ResponseError{common.ERR6, common.ResCodeDict[common.ERR6]}
			return common.RespondError(resErr)
		}
		Adminlist = append(Adminlist, admin)
	}

	adminJson, err2 := json.Marshal(Adminlist)
	if err2 != nil {
		//convert JSON eror
		resErr := common.ResponseError{common.ERR6, common.ResCodeDict[common.ERR6]}
		return common.RespondError(resErr)
	}

	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(adminJson)}
	return common.RespondSuccess(resSuc)
}

// ------------------- //
