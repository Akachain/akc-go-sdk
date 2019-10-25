module github.com/Akachain/akc-go-sdk

go 1.12

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b

require (
	github.com/Akachain/akc-go-sdk/common v0.0.0-20191025071001-7cf46a80ab48
	github.com/Akachain/akc-go-sdk/util v0.0.0-20191025071001-7cf46a80ab48
	github.com/hyperledger/fabric v1.4.3
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/tedsuo/ifrit v0.0.0-20180802180643-bea94bb476cc // indirect
	golang.org/x/tools/gopls v0.1.3 // indirect
)
