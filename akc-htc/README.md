
# Akachain High Throughput Chaincode Template

## I. Overview

The Akachain High Throughput Chaincode (AKC HTC) is designed for applications handling hundreds or thousands transaction per second which all read or update the same asset (key) in the ledger.

This document provides the AKC HTC template interface and how to use.

## II. AKC HTC Interface

The AKC HTC template is packaged into akc_htc package which provide the following interfaces

##### Insert: The insert function inserts the value into the temporary storage (the state db that may be deleted later) as a single row. The key is unique and created by combining the input and transaction id.

```go
Insert(<name>, <key>, <value>, <operation>)
```

- Name: The name, object or attribute that applied the high throughput chaincode. Example: merchant, user ...
- Key: The key identify the object. Example: merchant address, user id ...
- Operation: The operation that used for aggregation. Currently support OP_ADD (+) and OP_SUB (-)
- Value: The value of key. Currently for aggregation purpose only, so it should be in numeric type

#### Get: Get the value from temporary storage

```go
Get(<name>, [key])
```

-   Name: same as insert function
-   Key: [optional] same as insert function. If the key is null, all key:value associated with the name will be returned
##### Prune: Prune the temporary storage by aggregating the multiple row (value) into single row (value)

```go
Prune(<name>, [key], [prunt_type])
```

- Name: same as insert function
- Key: same as insert function. If key is null, all key:value associated with the name will be pruned
- Prunt_type: Currently, we support two type of prune:
	+ PRUNE_FAST: perform the aggregation operation then delete the related row
	+ PRUNE_SAFE: Same to PRUNE_FAST but the result is backup before delete all related row.

#### Delete: Delete the temporary storage

```go
Delete(<name>, [key])
```

- Name: same as insert function
- Key: same as insert function. If key is null, all key:value associated with the name will be deleted


## III. How to use

To use AKC HTC, the package akchtc must be imported to chaincode file. Ex:

```go
import (
  "encoding/json"
  "errors"
  "fmt"
  "reflect"
  "github.com/hyperledger/fabric/core/chaincode/shim"
  akchtc "github.com/Akachain/akc-go-sdk/akc-htc"
)

type ResponseData struct {
	Key  string
	Data []string
}

// example code insert using Akachain High throughput
func insertHTC(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  // Init Akachain High Throughput
  akc := akchtc.AkcHighThroughput{}
  res := akc.Insert(stub, []string{"variableName", "112233", "100", "OP_ADD"})

  if res != nil {
    return fmt.Sprintf("Failure"), res
  }
  return fmt.Sprintf("Success"), nil
}


// Example get data HTC with variableName
// This func response JSON data for name: "variableName"
func getHTC(stub shim.ChaincodeStubInterface, args []string) (string, error) {
  akc := akchtc.AkcHighThroughput{}
  res, err := akc.Get(stub, []string{"variableName"})

  var data map[string]ResponseData
	if err := json.Unmarshal(check, &data); err != nil {
		panic(err)
  }
  
  return fmt.Sprintf("%v", data), nil
}

// Example prune data HTC
// This func response JSON data for "variableName" after prune success.
func pruneHTC(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	akc := AkcHighThroughput{}
	resp, err := akc.Prune(stub, []string{"variableName"})

  if err != nil {
    return nil, err
  }

	return fmt.Sprintf("%s", resp), err
}

// Example delete variable "variableName" in HTC
func deleteHTC(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	akc := AkcHighThroughput{}
	resp, err := akc.Delete(stub, []string{"variableName"})

  if err != nil {
    return nil, err
  }

	return fmt.Sprintf("%s", resp), err
}
```