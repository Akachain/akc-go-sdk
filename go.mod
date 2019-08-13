module github.com/Akachain/akc-go-sdk

go 1.12

replace github.com/Akachain/akc-go-sdk/common v0.0.0 => ./common

replace github.com/Akachain/akc-go-sdk/util v0.0.0 => ./util

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b

require (
	github.com/Akachain/akc-go-sdk/akchtc v0.0.0-20190801094203-7616438f5374 // indirect
	github.com/Akachain/akc-go-sdk/common v0.0.0
	github.com/Akachain/akc-go-sdk/util v0.0.0
	github.com/hyperledger/fabric v1.4.2
	github.com/mitchellh/mapstructure v1.1.2
	github.com/stretchr/testify v1.3.0
)
