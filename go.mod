module github.com/Akachain/akc-go-sdk

go 1.12

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b

replace akc-go-sdk/common v0.0.0 => ./common

replace akc-go-sdk/util v0.0.0 => ./util

require (
	akc-go-sdk/common v0.0.0
	akc-go-sdk/util v0.0.0
	github.com/Shopify/sarama v1.24.1 // indirect
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/hyperledger/fabric v1.4.3
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/spf13/cobra v0.0.5 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tedsuo/ifrit v0.0.0-20180802180643-bea94bb476cc // indirect
	github.com/willf/bitset v1.1.10 // indirect
	golang.org/x/tools/gopls v0.1.3 // indirect
)
