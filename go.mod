module github.com/bitsongofficial/go-bitsong

go 1.16

require (
	github.com/bitsongofficial/chainmodules v0.0.0-20210705143912-90636d0f2a23
	github.com/cosmos/cosmos-sdk v0.42.6
	github.com/cosmos/go-bip39 v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.11
	github.com/tendermint/tm-db v0.6.4
	github.com/tendermint/tmlibs v0.9.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
