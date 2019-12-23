module akc-go-sdk/util

go 1.13

replace akc-go-sdk/common v0.0.0 => ../common

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b

require (
	akc-go-sdk/common v0.0.0
	github.com/hyperledger/fabric v1.4.4
	github.com/hyperledger/fabric-lib-go v1.0.0 // indirect
	github.com/ijc/Gotty v0.0.0-20170406111628-a8b993ba6abd // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.6.1
)
