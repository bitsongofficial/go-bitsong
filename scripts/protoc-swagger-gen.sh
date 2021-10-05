#!/usr/bin/env bash

set -eo pipefail

SDK_VERSION=v0.42.10

chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/third_party/proto

rm -rf ./tmp-swagger-gen ./tmp && mkdir -p ./tmp-swagger-gen ./tmp/proto ./tmp/third_party

cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto ./tmp
cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/third_party/proto ./tmp/third_party

proto_dirs=$(find ./tmp/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do

  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 -name 'query.proto')
  if [[ $dir =~ "cosmos" ]]; then
    query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  fi
  if [[ ! -z "$query_file" ]]; then
    protoc  \
    -I "tmp/proto" \
    -I "tmp/third_party/proto" \
    "$query_file" \
    --swagger_out=./tmp-swagger-gen \
    --swagger_opt=logtostderr=true --swagger_opt=fqn_for_swagger_name=true --swagger_opt=simple_operation_ids=true
  fi
done

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./swagger/config.json -o ./swagger/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
rm -rf ./tmp-swagger-gen

# clean proto files
rm -rf ./tmp