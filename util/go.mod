module github.com/Akachain/akc-go-sdk/util

go 1.12

replace akc-go-sdk/common v0.0.0 => ../common

require (
	akc-go-sdk/common v0.0.0
	github.com/Knetic/govaluate v3.0.0+incompatible // indirect
	github.com/fsouza/go-dockerclient v1.4.2 // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/hyperledger/fabric v1.4.2
	github.com/hyperledger/fabric-amcl v0.0.0-20181230093703-5ccba6eab8d6 // indirect
	github.com/hyperledger/fabric-lib-go v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.4.0 // indirect
	github.com/sykesm/zap-logfmt v0.0.2 // indirect
	golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7 // indirect
	google.golang.org/grpc v1.23.0 // indirect
)
