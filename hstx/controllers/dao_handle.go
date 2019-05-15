package controllers

import (
	"encoding/json"
	"fmt"

	"gitlab.com/akachain/akc-go-sdk/common"
	"gitlab.com/akachain/akc-go-sdk/hstx/models"
	"gitlab.com/akachain/akc-go-sdk/util"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"
)

type Proposal models.Proposal
type Quorum models.Quorum
type Commit models.Commit
type Admin models.Admin

//High secure transaction database handle
// ------------------- //

//create data
func create_data_(stub shim.ChaincodeStubInterface, TableModel string, row_key []string, data interface{}) error {
	var old_data interface{}
	row_was_found, err := util.InsertTableRow(stub, TableModel, row_key, data, util.FAIL_BEFORE_OVERWRITE, &old_data)
	if err != nil {
		return err
	}
	if row_was_found {
		return fmt.Errorf("Could not create data %v because an data already exists", data)
	}
	return nil //success
}

//get information of data  by ID
func get_data_byid_(stub shim.ChaincodeStubInterface, ID string, MODELTABLE string) (interface{}, error) {
	var datastruct interface{}

	row_was_found, err := util.GetTableRow(stub, MODELTABLE, []string{ID}, &datastruct, util.FAIL_IF_MISSING)

	if err != nil {
		return nil, err
	}
	if !row_was_found {
		return nil, fmt.Errorf("Data with ID %s does not exist", ID)
	}
	return datastruct, nil
}

//get all data
func get_all_data_(stub shim.ChaincodeStubInterface, MODELTABLE string) (chan []byte, error) {
	row_json_bytes, err := util.GetTableRows(stub, MODELTABLE, []string{})
	if err != nil {
		return nil, fmt.Errorf("Could not get %v", err.Error())
	}
	return row_json_bytes, nil
}

// GetDataByID
func GetDataByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	DataID := args[0]
	var data interface{}
	ModelTable := args[1]

	rs, err := get_data_byid_(stub, DataID, ModelTable)

	switch ModelTable {
	case models.ADMINTABLE:
		data = new(Admin)
	case models.COMMITTABLE:
		data = new(Commit)
	case models.QUORUMTABLE:
		data = new(Quorum)
	case models.PROPOSALTABLE:
		data = new(Proposal)
	default:
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	mapstructure.Decode(rs, data)
	fmt.Printf("Proposal: %v\n", data)

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
func GetAllData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}
	var data interface{}
	var Datalist []interface{}
	ModelTable := args[0]

	switch ModelTable {
	case models.ADMINTABLE:
		data = new(Admin)
	case models.COMMITTABLE:
		data = new(Commit)
	case models.QUORUMTABLE:
		data = new(Quorum)
	case models.PROPOSALTABLE:
		data = new(Proposal)
	default:
		//Invalid arguments
		resErr := common.ResponseError{common.ERR2, common.ResCodeDict[common.ERR2]}
		return common.RespondError(resErr)
	}

	datalbytes, err := get_all_data_(stub, ModelTable)
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
