package controllers

import (
	"fmt"

	"gitlab.com/akachain/akc-go-sdk/hstx/models"
	"gitlab.com/akachain/akc-go-sdk/util"

	"github.com/hyperledger/fabric/core/chaincode/shim"
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

// ------------------- //
