package pvtdata

import (
	"encoding/json"
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gotest.tools/assert"
	"os"
	"strconv"
	"strings"
	"testing"
)

//var stub *util.MockStubExtend
var collectionMarbles = "collectionMarbles"

type collectionsConfig struct {
	Name              string `json:"name"`
	Policy            string `json:"policy"`
	RequiredPeerCount int    `json:"requiredPeerCount"`
	MaxPeerCount      int    `json:"maxPeerCount"`
	BlockToLive       int    `json:"blockToLive"`
	MemberOnlyRead    bool   `json:"memberOnlyRead"`
}

type dbConfig struct {
	DbURL string `json:"TEST_COUCHDB_URL"`
}

func getCollectionsConfig(fileName string) ([]collectionsConfig, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	var ccfg []collectionsConfig
	// decode to get config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ccfg)
	if err != nil {
		return nil, err
	}

	return ccfg, nil
}

func getDBConfig(fileName string) (dbConfig, error) {
	// fileName is the path to the json config file
	file, err := os.Open(fileName)
	if err != nil {
		return dbConfig{}, err
	}

	var cfg dbConfig
	// decode to get config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

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

func TestDeletePrivateData(t *testing.T) {
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

	err = stub.DelPrivateData(collectionMarbles, marblesPvtInput.Name)
	assert.NilError(t, err)

	// Get private data information
	res, err := stub.GetPrivateData(collectionMarbles, marblesPvtInput.Name)
	assert.Check(t, err == nil)
	t.Log(string(res))
}

func TestGetPrivateDataByRange(t *testing.T) {
	stub, err := setupMockStub()
	assert.NilError(t, err)

	// Initialize data
	for i := 0; i < 10; i++ {
		name := "marbles" + strconv.Itoa(i)
		// Initialize args
		marblesInput := &marble{
			ObjectType: "marble",
			Name:       name,
			Color:      "blue",
			Size:       35 + i,
			Owner:      "blob",
		}

		// Initialize private data
		transient := make(map[string][]byte, 0)
		marblesPvtInput := &marblePrivateDetails{
			ObjectType: "marblePrivatePrice",
			Name:       name,
			Price:      1991 + i,
		}
		marblesBytes, err := json.Marshal(marblesPvtInput)
		assert.NilError(t, err)

		transient["marble"] = marblesBytes

		// Execution invoke method initMarble
		out := util.MockInvokePrivateTransaction(t, stub, [][]byte{[]byte("initMarble"), []byte(marblesInput.Name),
			[]byte(marblesInput.Color), []byte(strconv.Itoa(marblesInput.Size)), []byte(marblesInput.Owner)}, transient)
		assert.Equal(t, "", out)
	}

	// Query private data from marbles3 to marbles 8
	startKey := "marbles3"
	endKey := "marbles8"
	resIterator, err := stub.GetPrivateDataByRange(collectionMarbles, startKey, endKey)
	assert.NilError(t, err)

	i := 0
	for resIterator.HasNext() {
		resIterator.Next()
		i++
	}
	assert.Equal(t, i, 5)
}

func TestGetPrivateDataQueryResult(t *testing.T) {
	stub, err := setupMockStub()
	assert.NilError(t, err)

	// Initialize data
	for i := 0; i < 10; i++ {
		name := "marbles" + strconv.Itoa(i)
		// Initialize args
		marblesInput := &marble{
			ObjectType: "marble",
			Name:       name,
			Color:      "blue",
			Size:       35 + i,
			Owner:      "blob",
		}

		// Initialize private data
		transient := make(map[string][]byte, 0)
		marblesPvtInput := &marblePrivateDetails{
			ObjectType: "marblePrivatePrice",
			Name:       name,
			Price:      1991 + i,
		}
		marblesBytes, err := json.Marshal(marblesPvtInput)
		assert.NilError(t, err)

		transient["marble"] = marblesBytes

		// Execution invoke method initMarble
		out := util.MockInvokePrivateTransaction(t, stub, [][]byte{[]byte("initMarble"), []byte(marblesInput.Name),
			[]byte(marblesInput.Color), []byte(strconv.Itoa(marblesInput.Size)), []byte(marblesInput.Owner)}, transient)
		assert.Equal(t, "", out)
	}

	// Query private data that have price great than 1996
	queryString := `{"selector":{"price": {"$gt": 1996}}}`
	resIterator, err := stub.GetPrivateDataQueryResult(collectionMarbles, queryString)
	assert.NilError(t, err)

	i := 0
	for resIterator.HasNext() {
		resIterator.Next()
		i++
	}
	assert.Equal(t, i, 4)
}
