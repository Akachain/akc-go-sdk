# akc-go-sdk

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
    "github.com/hyperledger/fabric/core/chaincode/shim"
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

##4. Private Data Collection (pvtdata)
First of all we need define configuration for each private data collection. Example we have two private data collections are: "collectionMarbles" and "collectionMarblePrivateDetails"
```
[
 	{
	 	"name": "collectionMarbles",
	 	"policy": "OR('Org1MSP.member', 'Org2MSP.member')",
	 	"requiredPeerCount": 1,
	 	"maxPeerCount": 2,
	 	"blockToLive":1000000,
		"memberOnlyRead": false
	},
 	{
	 	"name": "collectionMarblePrivateDetails",
	 	"policy": "OR('Org2MSP.member', 'Org3MSP.member')",
	 	"requiredPeerCount": 1,
	 	"maxPeerCount": 2,
	 	"blockToLive":1000000,
		"memberOnlyRead": false
 	}
]
```
For Unit testing chaincode using private data collection API we need to add a list of imports that are neccessary as follows.

```
import (
    ...
    "testing"
    "github.com/Akachain/akc-go-sdk/util"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/stretchr/testify/assert"
)
```

We then create a *MockStubExtend* object that literally extend Fabric MockStub and chaincode:
```
func setupMockStub() (*util.MockStubExtend, error) {
	// Fetch database configuration
	dbCfg, err := getDBConfig("./testdata/config.json")
	if err != nil {
		return nil, err
	}

	// Fetch private data collections configuration
	ccfg, err := getCollectionsConfig("./testdata/collections_config.json")
	if err != nil {
		return nil, err
	}

	mapDBHandlers := make(map[string]*util.CouchDBHandler, 0)
	// Initialize private database handler
	for _, cfg := range ccfg {
		db, err := util.NewCouchDBHandlerWithConnection(strings.ToLower(cfg.Name), true, dbCfg.DbURL)
		if err != nil {
			return nil, err
		}
		mapDBHandlers[cfg.Name] = db
	}

	// Initialize mockstub private extend
	cc := new(SimpleChaincode)
	stub := util.NewMockStubPrivateDataExtend(shim.NewMockStub("sample", cc), cc, mapDBHandlers)

	return stub, nil
}
```

Then we can perform Unit test for each chaincode invoke function normally. Here is an example of testing an invoke function:
```
func TestInitMarble(t *testing.T) {
	stub, err := setupMockStub()
	assert.NilError(t, err)

	// Initialize args
	marblesInput := &marble{
		ObjectType: "marble",
		Name:       "marbles1",
		Color:      "blue",
		Size:       35,
		Owner:      "blob",
	}

	// Initialize private data
	transient := make(map[string][]byte, 0)
	marblesPvtInput := &marblePrivateDetails{
		ObjectType: "marblePrivatePrice",
		Name:       "marbles1",
		Price:      1991,
	}
	marblesBytes, err := json.Marshal(marblesPvtInput)
	assert.NilError(t, err)

	transient["marble"] = marblesBytes

	// Execution invoke method initMarble
	out := util.MockInvokePrivateTransaction(t, stub, [][]byte{[]byte("initMarble"), []byte(marblesInput.Name),
		[]byte(marblesInput.Color), []byte(strconv.Itoa(marblesInput.Size)), []byte(marblesInput.Owner)}, transient)
	assert.Equal(t, "", out)

	// Get Public information of marbles
	res, err := stub.GetState(marblesInput.Name)
	assert.NilError(t, err)

	marblePublicChecking := &marble{}
	err = json.Unmarshal(res, marblePublicChecking)
	assert.NilError(t, err)
	assert.DeepEqual(t, marblesInput, marblePublicChecking)

	// Get private data information
	res, err = stub.GetPrivateData(collectionMarbles, marblesInput.Name)
	assert.NilError(t, err)

	marblePvtChecking := &marblePrivateDetails{}
	err = json.Unmarshal(res, marblePvtChecking)
	assert.NilError(t, err)
	assert.DeepEqual(t, marblesPvtInput, marblePvtChecking)
}
```