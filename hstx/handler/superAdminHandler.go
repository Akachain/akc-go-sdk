package handler

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/hstx/model"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"
)

// SuperAdminHanler ...
type SuperAdminHanler struct{}

// CreateSuperAdmin ...
func (sah *SuperAdminHanler) CreateSuperAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 3)

	superAdmin := new(model.SuperAdmin)
	err := json.Unmarshal([]byte(args[0]), superAdmin)
	if err != nil {
		// Return error: can't unmashal json
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR3,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine()),
		})
	}

	superAdmin.Status = "Active"

	common.Logger.Infof("Create SuperAdmin: %+v\n", superAdmin)
	err = util.Createdata(stub, model.SuperAdminTable, []string{superAdmin.SuperAdminID}, &superAdmin)
	if err != nil {
		resErr := common.ResponseError{
			ResCode: common.ERR5,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine()),
		}
		return common.RespondError(resErr)
	}

	bytes, err := json.Marshal(superAdmin)
	if err != nil {
		// Return error: can't mashal json
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR5,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine()),
		})
	}

	resSuc := common.ResponseSuccess{
		ResCode: common.SUCCESS,
		Msg:     common.ResCodeDict[common.SUCCESS],
		Payload: string(bytes)}
	return common.RespondSuccess(resSuc)
}

// GetAllSuperAdmin ...
func (sah *SuperAdminHanler) GetAllSuperAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	res := util.GetAllData(stub, new(model.SuperAdmin), model.SuperAdminTable)
	return res
}

// GetSuperAdminByID ...
func (sah *SuperAdminHanler) GetSuperAdminByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 1)

	superAdminID := args[0]
	res := util.GetDataByID(stub, superAdminID, new(model.SuperAdmin), model.SuperAdminTable)
	return res
}

//UpdateSuperAdmin ...
func (sah *SuperAdminHanler) UpdateSuperAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 1)

	tmpSuperAdmin := new(model.SuperAdmin)
	err := json.Unmarshal([]byte(args[0]), tmpSuperAdmin)
	if err != nil {
		// Return error: can't unmashal json
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR3,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine()),
		})
	}

	if len(tmpSuperAdmin.SuperAdminID) == 0 {
		resErr := common.ResponseError{
			ResCode: common.ERR13,
			Msg:     fmt.Sprintf("%s %s", common.ResCodeDict[common.ERR13], err.Error()),
		}
		return common.RespondError(resErr)
	}

	//get superAdmin information
	rawSuperAdmin, err := util.Getdatabyid(stub, tmpSuperAdmin.SuperAdminID, model.SuperAdminTable)
	if err != nil {
		resErr := common.ResponseError{
			ResCode: common.ERR4,
			Msg:     fmt.Sprintf("%s %s", common.ResCodeDict[common.ERR4], err.Error()),
		}
		return common.RespondError(resErr)
	}

	superAdmin := new(model.SuperAdmin)
	mapstructure.Decode(rawSuperAdmin, superAdmin)

	tmpSuperAdminVal := reflect.ValueOf(tmpSuperAdmin).Elem()
	superAdminVal := reflect.ValueOf(superAdmin).Elem()
	for i := 0; i < tmpSuperAdminVal.NumField(); i++ {
		fieldName := tmpSuperAdminVal.Type().Field(i).Name
		if len(tmpSuperAdminVal.Field(i).String()) > 0 {
			field := superAdminVal.FieldByName(fieldName)
			if field.CanSet() {
				field.SetString(tmpSuperAdminVal.Field(i).String())
			}
		}
	}

	err = util.Changeinfo(stub, model.SuperAdminTable, []string{superAdmin.SuperAdminID}, superAdmin)
	if err != nil {
		//Overwrite fail
		resErr := common.ResponseError{
			ResCode: common.ERR5,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine()),
		}
		return common.RespondError(resErr)
	}

	bytes, err := json.Marshal(superAdmin)
	if err != nil {
		// Return error: can't mashal json
		return common.RespondError(common.ResponseError{
			ResCode: common.ERR5,
			Msg:     fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR5], err.Error(), common.GetLine()),
		})
	}

	resSuc := common.ResponseSuccess{
		ResCode: common.SUCCESS,
		Msg:     common.ResCodeDict[common.SUCCESS],
		Payload: string(bytes)}
	return common.RespondSuccess(resSuc)
}
