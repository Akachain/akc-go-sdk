# Akachain Golang Software Development Kit

[![Go Report Card](https://goreportcard.com/badge/github.com/Akachain/akc-go-sdk)](https://goreportcard.com/report/github.com/Akachain/akc-go-sdk)

golang SDK that supports writing chaincodes on Akachain platform. In release v1.0, we introduce 3 different Software Development Kits (SDKs)

## 1. Unit testing framework

Hyperledger Fabric supports writing unit test for chaincode using GoMock. However, the default Mock file provided by Fabric does not support testing with CouchDB. Instead, key-value data is stored in an in-memory map. This does not allows developers to perform unit test on any function that relies on couchdb. 

**akc-go-sdk** provides utilities that override the default mock stub class to allow writing chaincode unit test that uses a local (or remote) couchdb instance.  

We first need to install Apache couchdb  at http://couchdb.apache.org/ 

There is also a list of imports that are neccessary as follows.

```
import (
    ...
    "testing"
    "github.com/Akachain/akc-go-sdk/util"
    "github.com/hyperledger/fabric-chaincode-go/shim"
    "github.com/stretchr/testify/assert"
)
```

We then create a *MockStubExtend* object that literally extend Fabric MockStub and 

```
type Chaincode struct {
}

func setupMock() *util.MockStubExtend {
    // Initialize mockstubextend
    cc := new(Chaincode)
    stub := util.NewMockStubExtend(shim.NewMockStub("mockstubextend", cc), cc)

    // Create a new database, Drop old database
    db, _ := util.NewCouchDBHandler("dbtest", true)
    stub.SetCouchDBConfiguration(db)
    return stub
}
```

Then we can perform Unit test for each chaincode invoke function normally. Here is an example of testing an invoke function:

```
func TestSample(t *testing.T) {
	stub := setupMock()

	usr := models.UserWallet{
		UserId:      "id1",
		PublicKey:   "pubkey",
	}

	param, _ := json.Marshal(usr)

	// Create a new user
	rs := util.MockInvokeTransaction(t, stub, [][]byte{[]byte("createUser"), []byte(param), []byte("2")})

	// Make a composite key that is similar with the one in couchdb
	compositeKey, _ := stub.CreateCompositeKey(models.USERWALLET_TABLE, []string{rs})

	// Check if the created user exist
	state, _ := stub.GetState(compositeKey)
	var ad models.UserWallet
	json.Unmarshal([]byte(state), &ad)

	// Check if the created user information is correct
	assert.Equal(t, usr.UserId, ad.UserId)
	assert.Equal(t, usr.PublicKey, ad.PublicKey)
}
```

## 2. High Throughput Chaincode (HTC)
Please follow the instruction [here](https://docs.google.com/document/d/18IpdA-Io7hLNZs7cjHig-6bp4dCt0F-sK1cF1pC_euw/edit?usp=sharing)

## 3. High Secure Transaction Chaincode (HSTX)
A high level description of HSTX is described [here](https://drive.google.com/open?id=1FDVoU8L2a2U8rISxWei_ITMx4JchrS6n). It is still in R&D phase.

We have a live demo [here](https://akc-sdk.akachain.io/#/home)