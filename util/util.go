package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect" // This is only used in InterfaceIsNil
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	uuid "github.com/satori/go.uuid"
)

// Taken from https://stackoverflow.com/questions/13901819/quick-way-to-detect-empty-values-via-reflection-in-go
func InterfaceIsNilOrIsZeroOfUnderlyingType(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func CheckChaincodeFunctionCallWellFormedness(args []string, expected_arg_count int) error {
	if len(args) == expected_arg_count {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Expected %v arguments, but got %v; args: %v", expected_arg_count, len(args), args))
	}
}

// Pass through to Sprintf
func MakeErrorRetval(error_message string, args ...interface{}) ([]byte, error) {
	return nil, fmt.Errorf(error_message, args...)
}

// This is effectively a strongly typed enum declaration.
type GetTableRow_FailureOption bool

const (
	DONT_FAIL_IF_MISSING GetTableRow_FailureOption = false
	FAIL_IF_MISSING      GetTableRow_FailureOption = true
)

// Implementation of GetTableKey that returns the composite key, if the row was found, and error.
func getTableRowAndCompositeKey(
	stub shim.ChaincodeStubInterface,
	table_name string,
	row_keys []string,
	row_value interface{},
	failure_option GetTableRow_FailureOption,
) (composite_key string, rowWasFound bool, err error) {
	// Initialize this to default not-found.
	rowWasFound = false

	// Form the composite key that will index this table row in the ledger state key/value store.
	composite_key, err = stub.CreateCompositeKey(table_name, row_keys)
	if err != nil {
		composite_key = ""
		err = fmt.Errorf("GetTableRow failed because stub.CreateCompositeKey failed with error %v", err)
		return
	}

	//     fmt.Printf("getTableRowAndCompositeKey; table_name = \"%s\", composite_key (may contain unprintable chars) = \"%s\", row_value = %v, InterfaceIsNilOrIsZeroOfUnderlyingType(row_value) = %v\n", table_name, composite_key, row_value, InterfaceIsNilOrIsZeroOfUnderlyingType(row_value))

	var bytes []byte
	bytes, err = stub.GetState(composite_key)
	if err != nil {
		// Regardless of failure option, we will be returning due to this error.
		if failure_option == FAIL_IF_MISSING {
			err = fmt.Errorf("GetTableRow failed because stub.GetState(%v) failed with error %v", composite_key, err)
		} else {
			err = nil
		}
		return
	}
	if bytes == nil {
		// Regardless of failure option, we will be returning due to this bytes == nil condition.
		if failure_option == FAIL_IF_MISSING {
			err = fmt.Errorf("GetTableRow failed because row with keys %v does not exist", row_keys)
		} else {
			err = nil
		}
		return
	}

	// If we got this far, then the row was found.
	rowWasFound = true

	// If row_value is not nil, attempt to unmarshal the row.
	if !InterfaceIsNilOrIsZeroOfUnderlyingType(row_value) {
		err = json.Unmarshal(bytes, row_value)
		if err != nil {
			err = fmt.Errorf("GetTableRow failed because json.Unmarshal failed with error %v", err)
			return
		}
	}

	// Return with success
	err = nil
	return
}

// If row_value is nil, then don't bother unmarshaling the data.  Thus a check for the
// presence of a particular table row can be done by specifying nil for row_value.
func GetTableRow(
	stub shim.ChaincodeStubInterface,
	table_name string,
	row_keys []string,
	row_value interface{},
	failure_option GetTableRow_FailureOption,
) (rowWasFound bool, err error) {
	_, rowWasFound, err = getTableRowAndCompositeKey(stub, table_name, row_keys, row_value, failure_option)
	return
}

// If row_value is nil, then don't bother unmarshaling the data.  Thus a check for the
// presence of a particular table row can be done by specifying nil for row_value.
func GetTableRows(
	stub shim.ChaincodeStubInterface,
	table_name string,
	row_keys []string,
) (chan []byte, error) {
	state_query_iterator, err := stub.GetStateByPartialCompositeKey(table_name, row_keys)
	if err != nil {
		return nil, fmt.Errorf("GetTableRow failed because stub.CreateCompositeKey failed with error %v", err)
	}

	rowJSONBytesChannel := make(chan []byte, 32) // TODO: 32 is arbitrary; is there some other reasonable buffer size?

	go func() {
		for state_query_iterator.HasNext() {
			query_result_kv, err := state_query_iterator.Next()
			if err != nil {
				panic("this should never happen probably")
			}
			rowJSONBytesChannel <- query_result_kv.Value
		}
		close(rowJSONBytesChannel)
	}()

	return rowJSONBytesChannel, nil
}

// This is effectively a strongly typed enum declaration.
type InsertTableRow_FailureOption uint8

const (
	DONT_FAIL_UPON_OVERWRITE InsertTableRow_FailureOption = 0
	FAIL_BEFORE_OVERWRITE    InsertTableRow_FailureOption = 1
	FAIL_UNLESS_OVERWRITE    InsertTableRow_FailureOption = 2
)

// NOTE: This is the current abstraction to port old v0.6 style tables to current non-tables style ledger.
// Note that row_value must be json.Marshal-able.
// If old_row_value is not nil and the requested row is present, then the row will be unmarshaled into
// old_row_value before the new value (specified by row_value).  Note that if FAIL_BEFORE_OVERWRITE
// is triggered, then old_row_value will contain the row that existed already that triggered the failure.
// If an error is returned, then the table will not have been modified (TODO: Probably need to verify this).
func InsertTableRow(
	stub shim.ChaincodeStubInterface,
	table_name string,
	row_keys []string,
	new_row_value interface{},
	failure_option InsertTableRow_FailureOption,
	old_row_value interface{},
) (rowWasFound bool, err error) {
	rowWasFound = false
	err = nil

	// Check that new_row_value is valid (must be specified)
	if InterfaceIsNilOrIsZeroOfUnderlyingType(new_row_value) {
		err = fmt.Errorf("InsertTableRow failed because new_row_value was nil")
		return
	}

	// Check for the row's presence and retrieve its value into old_row_value if specified
	composite_key, rowWasFound, err := getTableRowAndCompositeKey(stub, table_name, row_keys, old_row_value, DONT_FAIL_IF_MISSING)
	if err != nil {
		err = fmt.Errorf("InsertTableRow failed because getTableRowAndCompositeKey failed with error %v", err)
		return
	}

	// Process the failure_option
	if failure_option == FAIL_BEFORE_OVERWRITE && rowWasFound {
		err = fmt.Errorf("InsertTableRow failed because the row existed already and FAIL_BEFORE_OVERWRITE was specified")
		return
	} else if failure_option == FAIL_UNLESS_OVERWRITE && !rowWasFound {
		err = fmt.Errorf("InsertTableRow failed because the row did not yet exist and FAIL_UNLESS_OVERWRITE was specified")
		return
	}

	// Serialize Member struct as JSON
	bytes, err := json.Marshal(new_row_value)
	if err != nil {
		err = fmt.Errorf("InsertTableRow failed because json.Marshal failed with error %v", err)
		return
	}

	// Store the data in the ledger state
	err = stub.PutState(composite_key, bytes)
	if err != nil {
		err = fmt.Errorf("InsertTableRow failed because stub.PutState(%v) failed with error %v", composite_key, err)
		return
	}

	// Return with success.
	err = nil
	return
}

// UpdateTableRow is similar to InsertTableRow without re-checking if the row is already exist
func UpdateTableRow(
	stub shim.ChaincodeStubInterface,
	table_name string,
	row_keys []string,
	new_row_value interface{},
) (err error) {
	err = nil

	// Check that new_row_value is valid (must be specified)
	if InterfaceIsNilOrIsZeroOfUnderlyingType(new_row_value) {
		err = fmt.Errorf("InsertTableRow failed because new_row_value was nil")
		return
	}

	// Form the composite key that will index this table row in the ledger state key/value store.
	compositeKey, err := stub.CreateCompositeKey(table_name, row_keys)
	if err != nil {
		err = fmt.Errorf("GetTableRow failed because stub.CreateCompositeKey failed with error %v", err)
		return
	}

	// Serialize Member struct as JSON
	bytes, err := json.Marshal(new_row_value)
	if err != nil {
		err = fmt.Errorf("InsertTableRow failed because json.Marshal failed with error %v", err)
		return
	}

	// Store the data in the ledger state
	err = stub.PutState(compositeKey, bytes)
	if err != nil {
		err = fmt.Errorf("InsertTableRow failed because stub.PutState(%v) failed with error %v", compositeKey, err)
		return
	}

	// Return with success.
	err = nil
	return
}

// If old_row_value is not nil, then the table row will be unmarshaled into old_row_value before being deleted.
func DeleteTableRow(
	stub shim.ChaincodeStubInterface,
	table_name string,
	row_keys []string,
	old_row_value interface{},
	failure_option GetTableRow_FailureOption,
) (rowWasFound bool, err error) {
	rowWasFound = false
	err = nil

	// Check for the row's presence and retrieve its value into old_row_value if specified
	composite_key, rowWasFound, err := getTableRowAndCompositeKey(stub, table_name, row_keys, old_row_value, DONT_FAIL_IF_MISSING)
	if err != nil {
		err = fmt.Errorf("DeleteTableRow failed because getTableRowAndCompositeKey failed with error %v", err)
		return
	}

	// Process the failure_option
	if failure_option == FAIL_IF_MISSING && !rowWasFound {
		err = fmt.Errorf("DeleteTableRow failed because the row was not found and FAIL_IF_MISSING was specified")
		return
	}

	// Actually delete the row
	err = stub.DelState(composite_key)
	if err != nil {
		err = fmt.Errorf("DeleteTableRow failed because stub.DelState(%v) failed with error %v", composite_key, err)
		return
	}

	// Return with success
	err = nil
	return
}

// MockInvokeTransaction creates a mock invoke transaction using MockStubExtend
func MockInvokeTransaction(t *testing.T, stub *MockStubExtend, args [][]byte) string {
	txId := genTxID()
	res := stub.MockInvoke(txId, args)
	if res.Status != shim.OK {
		return string(res.Message)
	}
	// fmt.Println(res.Payload)
	return string(res.Payload)
}

// MockQueryTransaction creates a mock query transaction using MockStubExtend
func MockQueryTransaction(t *testing.T, stub *MockStubExtend, args [][]byte) string {
	txId := genTxID()
	res := stub.MockInvoke(txId, args)
	if res.Status != shim.OK {
		t.FailNow()
		return string(res.Message)
	}
	return string(res.Payload)
}

// MockIInit creates a mock invoke transaction using MockStubExtend
func MockInitTransaction(t *testing.T, stub *MockStubExtend, args [][]byte) string {
	txId := genTxID()
	res := stub.MockInit(txId, args)
	if res.Status != shim.OK {
		return string(res.Message)
	}
	return string(res.Payload)
}

// Generate random transaction ID
func genTxID() string {
	uid := uuid.Must(uuid.NewV4())
	txId := fmt.Sprintf("%s", uid)
	return txId
}
