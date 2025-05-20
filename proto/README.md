## Manual Protobuf build

set env variables:
```sh
GOLANG_PROTOBUF_VERSION=1.36.6
GRPC_GATEWAY_VERSION=1.16.0
```

run:
```sh
go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@v${GOLANG_PROTOBUF_VERSION}
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION}
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v${GRPC_GATEWAY_VERSION}
```
then run:
```sh
 git clone https://github.com/cosmos/gogoproto.git;
  cd gogoproto; \
  go mod download; \
  make install
```
once installed, in the root of your project
```sh
# requires buf to be installed: https://buf.build/docs/installation/
cd ../proto
buf mod update
cd ..
buf generate
```

then, move the generated proto file into the right places:
```sh
cp -r ./github.com/bitsongofficial/go-bitsong/x/* x/
# cp -r ./github.com/<any-other>/<proto-imports>/x/* x/
```

then you can clean the repo:
```sh
rm -rf ./github.com
```