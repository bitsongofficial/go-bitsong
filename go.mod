module github.com/bitsongofficial/go-bitsong

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.42.1
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/genproto v0.0.0-20210207032614-bba0dbe2a9ea // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
