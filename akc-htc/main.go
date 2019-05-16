package akchtc

import (
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

/**
 * Insert (updates) the ledger to include a new delta for a particular variable. If this is the first time
 * this variable is being added to the ledger, then its initial value is assumed to be 0. The arguments
 * to give in the args array are as follows:
 *	- args[0] -> name of the variable
 *	- args[1] -> key of the variable
 *	- args[2] -> new delta (float)
 *	- args[3] -> operation (currently supported are addition "+" and subtraction "-")
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
	if op != "+" && op != "-" {
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
func (akcStub *AkcHighThroughput) Get(APIstub shim.ChaincodeStubInterface, args []string) (float64, error) {
	// Check we have a valid number of args
	if len(args) != 2 {
		return 0, fmt.Errorf("Incorrect number of arguments, expecting 2")
	}

	name := args[0]
	key := args[1]
	// Get all deltas for the variable
	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("varName~key~op~value~txID", []string{name, key})
	if deltaErr != nil {
		return 0, fmt.Errorf(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return 0, fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Iterate through result set and compute final value
	var finalVal float64
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return 0, fmt.Errorf(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return 0, fmt.Errorf(splitKeyErr.Error())
		}

		// Retrieve the delta value and operation
		operation := keyParts[2]
		valueStr := keyParts[3]

		// Convert the value string and perform the operation
		value, convErr := strconv.ParseFloat(valueStr, 64)
		if convErr != nil {
			return 0, fmt.Errorf(convErr.Error())
		}

		switch operation {
		case "+":
			finalVal += value
		case "-":
			finalVal -= value
		default:
			return 0, fmt.Errorf(fmt.Sprintf("Unrecognized operation %s", operation))
		}
	}

	// return finalValue AKC-HP []byte(strconv.FormatFloat(finalVal, 'f', -1, 64))
	return finalVal, nil
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
func (akcStub *AkcHighThroughput) Prune(APIstub shim.ChaincodeStubInterface, args []string) (bool, error) {
	// Check we have a valid number of ars
	if len(args) != 3 {
		return false, fmt.Errorf("Incorrect number of arguments, expecting 3")
	}

	// Retrieve the name of the variable to prune
	name := args[0]
	key := args[1]
	pruneType := args[2]

	// Get all delta rows for the variable
	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("varName~key~op~value~txID", []string{name, key})

	if deltaErr != nil {
		return false, fmt.Errorf(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	if pruneType == "PRUNE_FAST" {
		// Check the variable existed
		if !deltaResultsIterator.HasNext() {
			return false, fmt.Errorf(fmt.Sprintf("No variable by the name %s exists", name))
		}

		// Iterate through result set computing final value while iterating and deleting each key
		var finalVal float64
		var i int
		for i = 0; deltaResultsIterator.HasNext(); i++ {
			// Get the next row
			responseRange, nextErr := deltaResultsIterator.Next()
			if nextErr != nil {
				return false, fmt.Errorf(nextErr.Error())
			}

			// Split the key into its composite parts
			_, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
			if splitKeyErr != nil {
				return false, fmt.Errorf(splitKeyErr.Error())
			}

			// Retrieve the operation and value
			operation := keyParts[2]
			valueStr := keyParts[3]

			// Convert the value to a float
			value, convErr := strconv.ParseFloat(valueStr, 64)
			if convErr != nil {
				return false, fmt.Errorf(convErr.Error())
			}

			// Delete the row from the ledger
			deltaRowDelErr := APIstub.DelState(responseRange.Key)
			if deltaRowDelErr != nil {
				return false, fmt.Errorf(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
			}

			// Add the value of the deleted row to the final aggregate
			switch operation {
			case "+":
				finalVal += value
			case "-":
				finalVal -= value
			default:
				return false, fmt.Errorf(fmt.Sprintf("Unrecognized operation %s", operation))
			}
		}

		// Update the ledger with the final value and return
		updateResp := akcStub.Insert(APIstub, []string{name, key, strconv.FormatFloat(finalVal, 'f', -1, 64), "+"})
		if updateResp != nil {
			return true, nil // return nil if prune success
		}

		return false, fmt.Errorf(fmt.Sprintf("Failed to prune variable: all rows deleted but could not update value to %f, variable no longer exists in ledger", finalVal))
	} else if pruneType == "PRUNE_SAFE" {
		// Get the var's value and process it
		getResp, err := akcStub.Get(APIstub, []string{name, key})
		if err != nil {
			return false, fmt.Errorf(fmt.Sprintf("Could not retrieve the value of %s before pruning, pruning aborted: %s", name, key))
		}

		valueStr := getResp
		// val, convErr := strconv.ParseFloat(getResp, 64)
		// if convErr != nil {
		// 	return false, fmt.Errorf(fmt.Sprintf("Could not convert the value of %s to a number before pruning, pruning aborted: %s", name, convErr.Error()))
		// }

		// Store the var's value temporarily
		backupPutErr := APIstub.PutState(fmt.Sprintf("%s_%s_PRUNE_BACKUP", name, key), f2barr(valueStr))
		if backupPutErr != nil {
			return false, fmt.Errorf(fmt.Sprintf("Could not backup the value of %s before pruning, pruning aborted: %s", name, backupPutErr.Error()))
		}

		// Delete each row
		var i int
		for i = 0; deltaResultsIterator.HasNext(); i++ {
			responseRange, nextErr := deltaResultsIterator.Next()
			if nextErr != nil {
				return false, fmt.Errorf(fmt.Sprintf("Could not retrieve next row for pruning: %s", nextErr.Error()))
			}

			deltaRowDelErr := APIstub.DelState(responseRange.Key)
			if deltaRowDelErr != nil {
				return false, fmt.Errorf(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
			}
		}

		// Insert new row for the final value
		vStr := fmt.Sprintf("%f", valueStr)

		updateResp := akcStub.Insert(APIstub, []string{name, key, vStr, "+"})
		if updateResp != nil {
			return false, fmt.Errorf(fmt.Sprintf("Could not insert the final value of the variable after pruning, variable backup is stored in %s_PRUNE_BACKUP: %s", name, key))
		}

		// Delete the backup value
		delErr := APIstub.DelState(fmt.Sprintf("%s_%s_PRUNE_BACKUP", name, key))
		if delErr != nil {
			return false, fmt.Errorf(fmt.Sprintf("Could not delete backup value %s_PRUNE_BACKUP, this does not affect the ledger but should be removed manually", name))
		}

		return true, nil
	} else {
		return false, fmt.Errorf("Incorect option for prune or something else ! Try again.")
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
