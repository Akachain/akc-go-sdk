/**
This package is built based on high throughput chaincode sample: https://github.com/hyperledger/fabric-samples/blob/release/high-throughput/chaincode/high-throughput.go
We have modified to be used as a shared package for chaincode golang in fabric.
New feature add to package:
	- Add key to correct identifier for variable name. Key is optional.
	- Changed operator to text.
	- Returns JSON data instead of text when Get, Prune by variables name.
	- Add optional for Prune.
	- Returns all data of variable name if not added Key when Get or Prune.
*/

package akchtc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type AkcHighThroughput struct {
	Name      string
	Key       string
	Value     string
	Operation string
}

type ResponseData struct {
	Key  string
	Data []string
}

/**
 * Insert (updates) the ledger to include a new delta for a particular variable. If this is the first time
 * this variable is being added to the ledger, then its initial value is assumed to be 0. The arguments
 * to give in the args array are as follows:
 *	- args[0] -> name of the variable
 *	- args[1] -> key of the variable
 *	- args[2] -> new delta (float)
 *	- args[3] -> operation (currently supported are addition "OP_ADD" and subtraction "OP_SUB")
 *
 * @param APIstub The chaincode shim
 * @param args The arguments array for the update invocation
 *
 * @return A response structure indicating success or failure with a message
 */
func (akcStub *AkcHighThroughput) Insert(APIstub shim.ChaincodeStubInterface, args []string) error {
	// Check we have a valid number of args
	if len(args) != 4 {
		return fmt.Errorf("Incorrect number of arguments, expecting 4")
	}

	// Extract the args
	name := args[0]
	key := args[1]
	value := args[2]
	op := args[3]

	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("Provided value was not a number: %s", err)
	}

	// Make sure a valid operator is provided
	if op != "OP_ADD" && op != "OP_SUB" {
		return fmt.Errorf(fmt.Sprintf("Operator %s is unrecognized", op))
	}

	// Retrieve info needed for the update procedure
	txid := APIstub.GetTxID()
	compositeIndexName := "varName~key~op~value~txID"

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey, compositeErr := APIstub.CreateCompositeKey(compositeIndexName, []string{name, key, op, value, txid})
	if compositeErr != nil {
		return fmt.Errorf(fmt.Sprintf("Could not create a composite key for %s: %s", name, compositeErr.Error()))
	}

	// Save the composite key index
	compositePutErr := APIstub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return fmt.Errorf(fmt.Sprintf("Could not put operation for %s in the ledger: %s", name, compositePutErr.Error()))
	}

	return nil // return error = nil if add composite success
}

/**
 * Retrieves the aggregate value of a variable in the ledger. Gets all delta rows for the variable
 * and computes the final value from all deltas. The args array for the invocation must contain the
 * following argument:
 *	- args[0] -> The name of the variable to get the value of
 *	- args[1] -> The key of the variable to get the value of
 *
 * @param APIstub The chaincode shim
 * @param args The arguments array for the get invocation
 *
 * @return A response structure indicating success or failure with a message
 */
func (akcStub *AkcHighThroughput) Get(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// Check we have a valid number of args
	if len(args) < 1 {
		return []byte(""), fmt.Errorf("Incorrect number of arguments, expecting 1")
	}

	name := args[0]
	var queryArgs []string

	queryArgs = append(queryArgs, name)
	if len(args) == 2 {
		queryArgs = append(queryArgs, args[1])
	}

	// Get all deltas for the variable
	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("varName~key~op~value~txID", queryArgs)
	if deltaErr != nil {
		return []byte(""), fmt.Errorf(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return []byte(""), fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Iterate through result set and compute final value
	var dataResult ResponseData
	var responseMap map[string]ResponseData
	responseMap = make(map[string]ResponseData)
	var i int

	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return []byte(""), fmt.Errorf(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return []byte(""), fmt.Errorf(splitKeyErr.Error())
		}

		// Retrieve the delta value and operation
		key := keyParts[1]
		operation := keyParts[2]
		valueStr := keyParts[3]
		value, convErr2 := strconv.ParseFloat(valueStr, 64)
		if convErr2 != nil {
			return []byte(""), fmt.Errorf(convErr2.Error())
		}

		if responseMap[key].Key == key {
			mapValue, convErr := strconv.ParseFloat(responseMap[key].Data[0], 64)
			if convErr != nil {
				return []byte(""), fmt.Errorf(convErr.Error())
			}

			returnValue, errVal := akcStub.getValueByOperator(operation, mapValue, value)

			if errVal != nil {
				return nil, errVal
			}

			responseMap[key].Data[0] = returnValue
		} else {
			firstValue, _ := strconv.ParseFloat("0", 64)
			dataResult.Key = key

			returnValue, errVal := akcStub.getValueByOperator(operation, firstValue, value)

			if errVal != nil {
				return nil, errVal
			}

			dataResult.Data = []string{returnValue}
			responseMap[key] = dataResult
		}
	}

	jsonRes, _ := json.Marshal(responseMap)

	// return finalValue AKC-HP []byte(strconv.FormatFloat(finalVal, 'f', -1, 64))
	return []byte(jsonRes), nil
}

/**
 * Prune a variable by deleting all of its delta rows while computing the final value. Once all rows
 * have been processed and deleted, a single new row is added which defines a delta containing the final
 * computed value of the variable. If type prune is PRUNE_FAST, this is NOT safe as any failures or errors during pruning
 * will result in an undefined final value for the variable and loss of data. Use type PRUNE_SAFE if data
 * integrity is important. The args array contains the following argument:
 *	- args[0] -> The name of the variable to prune
 *	- args[1] -> The key of the variable to prune
 *	- args[2] -> Type of prune
 *
 * @param APIstub The chaincode shim
 * @param args The args array for the pruneFast invocation
 *
 * @return A response structure indicating success or failure with a message
 */
func (akcStub *AkcHighThroughput) Prune(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// Check we have a valid number of ars
	if len(args) < 2 {
		return []byte(""), fmt.Errorf("Incorrect number of arguments, expecting 2")
	}

	// Retrieve the name of the variable to prune
	name := args[0]

	var key, pruneType string
	var queryStr []string

	queryStr = append(queryStr, name)
	if len(args) == 3 {
		key = args[1]
		pruneType = args[2]
		queryStr = append(queryStr, key)
	} else {
		pruneType = args[1]
	}

	if pruneType != "PRUNE_FAST" && pruneType != "PRUNE_SAFE" {
		return []byte(""), fmt.Errorf(fmt.Sprintf("Prune type %s is not supported", pruneType))
	}

	// Get all delta rows for the variable
	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("varName~key~op~value~txID", queryStr)

	if deltaErr != nil {
		return []byte(""), fmt.Errorf(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	if pruneType == "PRUNE_FAST" {
		// Check the variable existed
		if !deltaResultsIterator.HasNext() {
			return []byte(""), fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", name))
		}

		// Iterate through result set computing final value while iterating and deleting each key
		var dataResult ResponseData
		var responseMap map[string]ResponseData
		responseMap = make(map[string]ResponseData)
		var i int

		for i = 0; deltaResultsIterator.HasNext(); i++ {
			// Get the next row
			responseRange, nextErr := deltaResultsIterator.Next()
			if nextErr != nil {
				return []byte(""), fmt.Errorf(nextErr.Error())
			}

			// Split the key into its composite parts
			_, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
			if splitKeyErr != nil {
				return []byte(""), fmt.Errorf(splitKeyErr.Error())
			}

			// Retrieve the operation and value
			key = keyParts[1]
			operation := keyParts[2]
			valueStr := keyParts[3]

			if responseMap[key].Key == key {
				mapValue, convErr := strconv.ParseFloat(responseMap[key].Data[0], 64)
				if convErr != nil {
					return []byte(""), fmt.Errorf(convErr.Error())
				}

				value, convErr2 := strconv.ParseFloat(valueStr, 64)
				if convErr2 != nil {
					return []byte(""), fmt.Errorf(convErr2.Error())
				}

				switch operation {
				case "OP_ADD":
					mapValue += value
				case "OP_SUB":
					mapValue -= value
				default:
					return []byte(""), fmt.Errorf(fmt.Sprintf("Unrecognized operation %s", operation))
				}

				responseMap[key].Data[0] = fmt.Sprintf("%v", mapValue)
			} else {
				dataResult.Key = key
				dataResult.Data = []string{valueStr}
				responseMap[key] = dataResult
			}

			// Delete the row from the ledger
			deltaRowDelErr := APIstub.DelState(responseRange.Key)
			if deltaRowDelErr != nil {
				return []byte(""), fmt.Errorf(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
			}
		}

		// Loop and destroy
		for _, data := range responseMap {
			keyStrMap := data.Key
			valueStrMap := data.Data[0]

			// Update the ledger with the final value and return
			_, updateErr := akcStub.pruneFastUpdate(APIstub, name, keyStrMap, valueStrMap)
			if updateErr != nil {
				return []byte(""), updateErr
			}
		}

		jsonRes, _ := json.Marshal(responseMap)
		return []byte(jsonRes), nil // return nil if prune success
	} else if pruneType == "PRUNE_SAFE" {
		// Get the var's value and process it
		getResp, err := akcStub.Get(APIstub, queryStr)
		if err != nil {
			return []byte(""), fmt.Errorf(fmt.Sprintf("Could not retrieve the value of %s before pruning, pruning aborted: %s", name, key))
		}

		// Unmarshal get response
		var responseData map[string]ResponseData
		if err := json.Unmarshal(getResp, &responseData); err != nil {
			panic(err)
		}

		// Loop and destroy
		for _, data := range responseData {
			keyStrMap := data.Key
			valueStrMap := data.Data[0]

			// PutState before delete
			_, errBackup := akcStub.pruneSafeBackup(APIstub, name, keyStrMap, valueStrMap)
			if errBackup != nil {
				return []byte(""), errBackup
			}

			// Delete each row
			var i int
			for i = 0; deltaResultsIterator.HasNext(); i++ {
				responseRange, nextErr := deltaResultsIterator.Next()
				if nextErr != nil {
					return []byte(""), fmt.Errorf(fmt.Sprintf("Could not retrieve next row for pruning: %s", nextErr.Error()))
				}

				deltaRowDelErr := APIstub.DelState(responseRange.Key)
				if deltaRowDelErr != nil {
					return []byte(""), fmt.Errorf(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
				}
			}

			// DelState after delete row
			_, errDel := akcStub.pruneSafeUpdate(APIstub, name, keyStrMap, valueStrMap)
			if errDel != nil {
				return []byte(""), errDel
			}
		}

		return getResp, nil
	} else {
		return []byte(""), fmt.Errorf("Incorect option for prune or something else ! Try again.")
	}
}

/**
 * Deletes all rows associated with an aggregate variable from the ledger. The args array
 * contains the following argument:
 *	- args[0] -> The name of the variable to delete
 *	- args[1] -> The key of the variable to prune
 *
 * @param APIstub The chaincode shim
 * @param args The arguments array for the delete invocation
 *
 * @return A response structure indicating success or failure with a message
 */
func (akcStub *AkcHighThroughput) Delete(APIstub shim.ChaincodeStubInterface, args []string) (bool, error) {
	// Check there are a correct number of arguments
	if len(args) != 2 {
		return false, fmt.Errorf("Incorrect number of arguments, expecting 2")
	}

	// Retrieve the variable name
	name := args[0]
	key := args[1]

	// Delete all delta rows
	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("varName~key~op~value~txID", []string{name, key})
	if deltaErr != nil {
		return false, fmt.Errorf(fmt.Sprintf("Could not retrieve delta rows for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Ensure the variable exists
	if !deltaResultsIterator.HasNext() {
		return false, fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Iterate through result set and delete all indices
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return false, fmt.Errorf(fmt.Sprintf("Could not retrieve next delta row: %s", nextErr.Error()))
		}

		deltaRowDelErr := APIstub.DelState(responseRange.Key)
		if deltaRowDelErr != nil {
			return false, fmt.Errorf(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
		}
	}

	return true, nil
}

/**
 * Converts a float64 to a byte array
 *
 * @param f The float64 to convert
 *
 * @return The byte array representation
 */
func f2barr(f float64) []byte {
	str := strconv.FormatFloat(f, 'f', -1, 64)

	return []byte(str)
}

/** Private function **/

func (akcStub *AkcHighThroughput) pruneSafeBackup(APIstub shim.ChaincodeStubInterface, name, key, valueStr string) (bool, error) {
	value, _ := strconv.ParseFloat(valueStr, 64)
	// Store the var's value temporarily
	backupPutErr := APIstub.PutState(fmt.Sprintf("%s_%s_PRUNE_BACKUP", name, key), f2barr(value))
	if backupPutErr != nil {
		return false, fmt.Errorf(fmt.Sprintf("Could not backup the value of %s before pruning, pruning aborted: %s", name, backupPutErr.Error()))
	}

	return true, nil
}

func (akcStub *AkcHighThroughput) pruneSafeUpdate(APIstub shim.ChaincodeStubInterface, name, key, valueStr string) (bool, error) {
	updateResp := akcStub.Insert(APIstub, []string{name, key, valueStr, "OP_ADD"})
	if updateResp != nil {
		return false, fmt.Errorf(fmt.Sprintf("Could not insert the final value of the variable after pruning, variable backup is stored in %s_PRUNE_BACKUP: %s", name, key))
	}

	// Delete the backup value
	delErr := APIstub.DelState(fmt.Sprintf("%s_%s_PRUNE_BACKUP", name, key))
	if delErr != nil {
		return false, fmt.Errorf(fmt.Sprintf("Could not delete backup value %s_PRUNE_BACKUP, this does not affect the ledger but should be removed manually", name))
	}

	return true, nil
}

func (akcStub *AkcHighThroughput) pruneFastUpdate(APIstub shim.ChaincodeStubInterface, name, key, valueStr string) (bool, error) {
	updateResp := akcStub.Insert(APIstub, []string{name, key, valueStr, "OP_ADD"})
	if updateResp != nil {
		return false, fmt.Errorf(fmt.Sprintf("Could not insert the final value of the variable after pruning, name %s, key %s", name, key))
	}

	return true, nil
}

func (akcStub *AkcHighThroughput) getValueByOperator(operation string, currentValue, value float64) (string, error) {
	switch operation {
	case "OP_ADD":
		currentValue += value
	case "OP_SUB":
		currentValue -= value
	default:
		return "0", fmt.Errorf(fmt.Sprintf("Unrecognized operation %s", operation))
	}

	return fmt.Sprintf("%v", currentValue), nil
}
