package util

import (
	"encoding/json"
	"fmt"

	"github.com/Akachain/akc-go-sdk/common"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"
)

//High secure transaction database handle
// ------------------- //

// Update  Infomation
func Changeinfo(stub shim.ChaincodeStubInterface, TableModel string, row_key []string, data interface{}) error {
	_, err := InsertTableRow(stub, TableModel, row_key, data, FAIL_UNLESS_OVERWRITE, nil)
	return err
}

//create data
func Createdata(stub shim.ChaincodeStubInterface, TableModel string, row_key []string, data interface{}) error {
	var old_data interface{}
	row_was_found, err := InsertTableRow(stub, TableModel, row_key, data, FAIL_BEFORE_OVERWRITE, &old_data)
	if err != nil {
		return err
	}
	if row_was_found {
		return fmt.Errorf("Could not create data %v because an data already exists", data)
	}
	return nil //success
}

//get information of data  by ID
func Getdatabyid(stub shim.ChaincodeStubInterface, ID string, MODELTABLE string) (interface{}, error) {
	var datastruct interface{}

	row_was_found, err := GetTableRow(stub, MODELTABLE, []string{ID}, &datastruct, FAIL_IF_MISSING)

	if err != nil {
		return nil, err
	}
	if !row_was_found {
		return nil, fmt.Errorf("Data with ID %s does not exist", ID)
	}
	return datastruct, nil
}

//get information of data  by row keys
func Getdatabyrowkeys(stub shim.ChaincodeStubInterface, rowKeys []string, MODELTABLE string) (interface{}, error) {
	var datastruct interface{}

	row_was_found, err := GetTableRow(stub, MODELTABLE, rowKeys, &datastruct, FAIL_IF_MISSING)

	if err != nil {
		return nil, err
	}
	if !row_was_found {
		return nil, fmt.Errorf("Data with rowKeys %s does not exist", rowKeys)
	}
	return datastruct, nil
}

//get all data
func Getalldata(stub shim.ChaincodeStubInterface, MODELTABLE string) (chan []byte, error) {
	row_json_bytes, err := GetTableRows(stub, MODELTABLE, []string{})
	if err != nil {
		return nil, fmt.Errorf("Could not get %v", err.Error())
	}
	return row_json_bytes, nil
}

// GetDataByID
func GetDataByID(stub shim.ChaincodeStubInterface, DataID string, data interface{}, ModelTable string) pb.Response {

	rs, err := Getdatabyid(stub, DataID, ModelTable)
	if rs != nil {
		mapstructure.Decode(rs, data)
	} else {
		data = nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		//Convert Json Fail
		resErr := common.ResponseError{common.ERR3, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	fmt.Printf("Response: %s\n", string(bytes))
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(bytes)}
	return common.RespondSuccess(resSuc)
}

// GetDataByKey
func GetDataByRowKeys(stub shim.ChaincodeStubInterface, rowKeys []string, data interface{}, ModelTable string) pb.Response {

	rs, err := Getdatabyrowkeys(stub, rowKeys, ModelTable)

	mapstructure.Decode(rs, data)
	fmt.Printf("data: %v\n", data)

	bytes, err := json.Marshal(data)
	if err != nil {
		//Convert Json Fail
		resErr := common.ResponseError{common.ERR3, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR3], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
	fmt.Printf("Response: %s\n", string(bytes))
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(bytes)}
	return common.RespondSuccess(resSuc)
}

// GetAllData
func GetAllData(stub shim.ChaincodeStubInterface, data interface{}, ModelTable string) pb.Response {

	var Datalist []interface{}

	datalbytes, err := Getalldata(stub, ModelTable)
	for row_json_bytes := range datalbytes {
		err = json.Unmarshal(row_json_bytes, data)
		if err != nil {

			resErr := common.ResponseError{common.ERR6, common.ResCodeDict[common.ERR6]}
			return common.RespondError(resErr)
		}
		Datalist = append(Datalist, data)
	}
	if err != nil {
		//Get data eror
		resErr := common.ResponseError{common.ERR3, common.ResCodeDict[common.ERR3]}
		return common.RespondError(resErr)
	}
	dataJson, err2 := json.Marshal(Datalist)
	if err2 != nil {
		//convert JSON eror
		resErr := common.ResponseError{common.ERR6, common.ResCodeDict[common.ERR6]}
		return common.RespondError(resErr)
	}
	fmt.Printf("Response: %s\n", string(dataJson))
	resSuc := common.ResponseSuccess{common.SUCCESS, common.ResCodeDict[common.SUCCESS], string(dataJson)}
	return common.RespondSuccess(resSuc)
}

// ------------------- //
