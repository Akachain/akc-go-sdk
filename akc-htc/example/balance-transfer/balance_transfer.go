/*
 * This chaincode is an example of how to use Akachain SDK high throughput in golang.
 *
 * @author	Huan Le for Akachain
 * @created	18 Jun 2019
 */

package example

import (
	"encoding/json"
	"fmt"
	"strconv"

	. "github.com/Akachain/akc-go-sdk/common"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/mitchellh/mapstructure"

	akc "github.com/Akachain/akc-go-sdk/akc-htc"
)

var logger = shim.NewLogger("example_balance_transfer")

// Chaincode example simple Chaincode implementation
type Chaincode struct {
}

// USERTABLE - Table name in onchain
const USERTABLE = "User_"

// User - struct model
type User struct {
	WalletAddress string  `json:"WalletAddress"`
	Amount        float64 `json:"Amount"`
}

// AKC high throughput response data
type AkcResponseData struct {
	Key  string
	Data []string
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### example_balance_transfer Init ###########")

	return shim.Success(nil)

}

// Invoke chaincode
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### example_balance_transfer Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "createUser":
		return t.createUser(stub, args)

	case "transferToken":
		return t.transferToken(stub, args)

	case "updateUserBalance":
		return t.updateUserBalance(stub, args)

	case "getUserToken":
		return t.getUserToken(stub, args)
	}

	resErr := ResponseError{ERR5, "Unknown action, check the first argument"}
	logger.Errorf("Unknown action, check the first argument. Got: %v", args[0])
	return RespondError(resErr)
}

/**
	* Create user and save to database onchain. All attributes of User defined on the model.
	*
	*	- args[0] -> user walletAddress. Optional, if this argument is empty, wallet address = txID.
  *	- args[1] -> token Amount initialize.
  *
  * @param APIstub The chaincode shim
	* @param args The arguments array for the update invocation
	*
 	* @return a string of wallet address of User when success or message when error.
*/
func (t *Chaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 2)
	WalletAddress := args[0]
	if WalletAddress == "" {
		WalletAddress = stub.GetTxID()
	}
	Amount, _ := strconv.ParseFloat(args[1], 64)

	// Create new user document and save to database onchain.
	err := util.Createdata(stub, USERTABLE, []string{WalletAddress}, &User{WalletAddress: WalletAddress, Amount: Amount})
	if err != nil {
		resErr := ResponseError{ERR5, fmt.Sprintf("%s %s", ResCodeDict[ERR5], "")}
		return RespondError(resErr)
	} else {
		// Update user data infor to AKC high throughput
		akchtc := akc.AkcHighThroughput{}
		akchtc.Insert(stub, []string{"User", WalletAddress, args[1], "OP_ADD"})
	}
	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], WalletAddress}
	return RespondSuccess(resSuc)
}

/**
	* Sending tokens from account A to account B. However, this transfer will not be done immediately,
	* that it will push into high throughput. In this example, I use Akachain SDK: akc-highthroughput.
	*
	*	- args[0] -> wallet address send
	*	- args[1] -> wallet address receive
	*	- args[2] -> token amount will send.
  *
  * @param APIstub The chaincode shim
	* @param args The arguments array for the update invocation
	*
	* @return a string transaction ID response structure indicating success or failure with a message
*/
func (t *Chaincode) transferToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 3)

	walletSend := args[0]
	walletReceived := args[1]
	amountToken := args[2]

	// Perform the execution
	_, err := strconv.Atoi(amountToken)
	if err != nil {
		resErr := ResponseError{ERR8, fmt.Sprintf("%s %s", ResCodeDict[ERR8], "Invalid transaction amount, expecting a integer value")}
		return RespondError(resErr)
	}

	// Init AKC sdk high throughput
	akchtc := akc.AkcHighThroughput{}
	insertHTP := akchtc.Insert(stub, []string{"User", walletSend, amountToken, "OP_SUB"})

	if insertHTP != nil {
		// Insert High throughput failure
		resErr := ResponseError{ERR8, fmt.Sprintf("%s %s", ResCodeDict[ERR8], insertHTP.Error())}
		return RespondError(resErr)
	} else {
		akchtc.Insert(stub, []string{"User", walletReceived, amountToken, "OP_ADD"})
	}

	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], stub.GetTxID()}
	return RespondSuccess(resSuc)
}

/**
	* As mentioned above, that transfer will not be done immediately, it's pushed into
	* AKC high throughput. And this function will perform that balance transfer via the AKC
	* high throughput SDK.
	*
	*	- args[0] -> type of prune <"PRUNE_SAFE" or "PRUNE_FAST">. (To perform calculation, AKC
	* sdk will prune delta data, so need to know the prune type)
  *
  * @param APIstub The chaincode shim
	* @param args The arguments array for the update invocation
	*
	* @return a string transaction ID response structure indicating success or failure with a message.
*/
func (t *Chaincode) updateUserBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 1)

	pruneType := args[0]

	// Init AKC sdk high throughput
	akchtc := akc.AkcHighThroughput{}
	pruneBalance, akcErr := akchtc.Prune(stub, []string{"User", pruneType})

	if akcErr != nil {
		resErr := ResponseError{"AKC_ERR", akcErr.Error()}
		return RespondError(resErr)
	}

	// Unmarshal get response
	var responseData map[string]AkcResponseData
	fmt.Printf("responseData: %s", responseData)
	if err := json.Unmarshal(pruneBalance, &responseData); err != nil {
		resErr := ResponseError{"AKC_ERR", err.Error()}
		return RespondError(resErr)
	}

	for _, data := range responseData {
		walletAddress := data.Key
		tokenAmount, _ := strconv.ParseFloat(data.Data[0], 64)

		var userData User
		// Get user data in database by ID of document. In this, ID is wallet address of user.
		user_rs, err := util.Getdatabyid(stub, walletAddress, USERTABLE)
		if err != nil {
			resErr := ResponseError{ERR4, fmt.Sprintf("%s %s", ResCodeDict[ERR4], err.Error())}
			return RespondError(resErr)
		}

		mapstructure.Decode(user_rs, &userData)

		userData.Amount = tokenAmount

		// Overwrite user info.
		err = util.Changeinfo(stub, USERTABLE, []string{walletAddress}, &User{WalletAddress: walletAddress, Amount: userData.Amount})
		if err != nil {
			// Overwrite fail
			resErr := ResponseError{ERR5, fmt.Sprintf("%s %s", ResCodeDict[ERR5], err.Error())}
			return RespondError(resErr)
		}
	}

	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], stub.GetTxID()}
	return RespondSuccess(resSuc)
}

/**
	* Get token amount of user via wallet address.
	*
	*	- args[0] -> wallet address.
  *
  * @param APIstub The chaincode shim
	* @param args The arguments array for the update invocation
	*
	* @return a string token amount of user when success or failure with a message.
*/
func (t *Chaincode) getUserToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	util.CheckChaincodeFunctionCallWellFormedness(args, 1)

	walletAddress := args[0]
	var userData User

	// Get user data in database by ID of document. In this, ID is wallet address of user.
	user_rs, err := util.Getdatabyid(stub, walletAddress, USERTABLE)
	if err != nil {
		resErr := ResponseError{ERR4, fmt.Sprintf("%s %s", ResCodeDict[ERR4], err.Error())}
		return RespondError(resErr)
	}
	mapstructure.Decode(user_rs, &userData)
	resSuc := ResponseSuccess{SUCCESS, ResCodeDict[SUCCESS], fmt.Sprintf("%v", userData.Amount)}
	return RespondSuccess(resSuc)
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		logger.Errorf("Error starting chaincode: %s", err)
	}
}
