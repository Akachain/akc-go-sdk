package util

import (
	"encoding/json"
	"fmt"

	"github.com/Akachain/akc-go-sdk/common"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/mitchellh/mapstructure"
)

// Update  Infomation
func Changeinfo(stub shim.ChaincodeStubInterface, TableModel string, row_key []string, data interface{}) error {
	_, err := InsertTableRow(stub, TableModel, row_key, data, FAIL_UNLESS_OVERWRITE, nil)
	return err
}

// UpdateExistingData works similar to Changeinfo.
// However, it does not check if the document is already existed
// This is useful if we already query out the row before and do not want to query again.
func UpdateExistingData(stub shim.ChaincodeStubInterface, TableModel string, row_key []string, data interface{}) error {
	err := UpdateTableRow(stub, TableModel, row_key, data)
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
	if err != nil {
		//Get Data Fail
		resErr := common.ResponseError{common.ERR4, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
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
	if err != nil {
		//Get Data Fail
		resErr := common.ResponseError{common.ERR4, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}
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

// GetAllData
func GetAllData(stub shim.ChaincodeStubInterface, data interface{}, ModelTable string) pb.Response {

	// var Datalist []interface{}
	var Datalist = make([]map[string]interface{}, 0)

	datalbytes, err := Getalldata(stub, ModelTable)
	if err != nil {
		//Get Data Fail
		resErr := common.ResponseError{common.ERR4, fmt.Sprintf("%s %s %s", common.ResCodeDict[common.ERR4], err.Error(), common.GetLine())}
		return common.RespondError(resErr)
	}

	for row_json_bytes := range datalbytes {

		err := json.Unmarshal(row_json_bytes, data)
		if err != nil {
			resErr := common.ResponseError{common.ERR3, common.ResCodeDict[common.ERR6]}
			return common.RespondError(resErr)
		}

		bytes, err := json.Marshal(data)
		if err != nil {
			//convert JSON eror
			resErr := common.ResponseError{common.ERR3, common.ResCodeDict[common.ERR6]}
			return common.RespondError(resErr)
		}

		var temp map[string]interface{}
		err = json.Unmarshal(bytes, &temp)
		if err != nil {
			resErr := common.ResponseError{common.ERR3, common.ResCodeDict[common.ERR6]}
			return common.RespondError(resErr)
		}
		Datalist = append(Datalist, temp)
	}

	//if err != nil {
	//	//Get data eror
	//	resErr := common.ResponseError{common.ERR3, common.ResCodeDict[common.ERR3]}
	//	return common.RespondError(resErr)
	//}

	fmt.Printf("Datalist: %v\n", Datalist)
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
