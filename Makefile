#!/usr/bin/make -f

include scripts/makefiles/build.mk
include scripts/makefiles/e2e.mk
include scripts/makefiles/hl.mk
include scripts/makefiles/docker.mk
include contrib/devtools/Makefile

.DEFAULT_GOAL := help
help:
	@echo "Available top-level commands:"
	@echo ""
	@echo "Usage:"
	@echo "    make [command]"
	@echo ""
	@echo "  make build                 Build Bitsong node binary"
	@echo "  make install               Install Bitsong node binary"
	@echo "  make hl                    Show available docker commands (via Strangelove's Heighliner Tooling)"
	@echo "  make e2e                   Show available e2e commands"
	@echo ""
	@echo "Run 'make [subcommand]' to see the available commands for each subcommand."

APP_DIR = ./app
BINDIR ?= $(GOPATH)/bin

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
PACKAGES_UNITTEST=$(shell go list ./... | grep -v '/simulation' | grep -v '/cli_test')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
BUILDDIR ?= $(CURDIR)/build

TENDERMINT_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::') # grab everything after the space in "github.com/tendermint/tendermint v0.34.7"

DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

E2E_UPGRADE_VERSION := "v0.18.0"

GO_MODULE := $(shell cat go.mod | grep "module " | cut -d ' ' -f 2)
GO_VERSION := $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f 2)
GO_MAJOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
GO_MINIMUM_MAJOR_VERSION = $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f2 | cut -d'.' -f1)
GO_MINIMUM_MINOR_VERSION = $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f2 | cut -d'.' -f2)
# message to be printed if Go does not meet the minimum required version
GO_VERSION_ERR_MSG = "ERROR: Go version $(GO_MINIMUM_MAJOR_VERSION).$(GO_MINIMUM_MINOR_VERSION)+ is required"

export GO111MODULE = on

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

# process build tags
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=go-bitsong \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=bitsongd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
          -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TENDERMINT_VERSION)

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)' -trimpath

all: install tools lint

build: go.sum
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/bitsongd.exe ./cmd/bitsongd
else
	go build $(BUILD_FLAGS) -o bin/bitsongd ./cmd/bitsongd
endif

build-linux: go.sum
	go build $(BUILD_FLAGS)

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/bitsongd

#update-swagger-docs: statik
#	$(BINDIR)/statik -src=swagger/swagger-ui -dest=swagger -f -m
#	@if [ -n "$(git status --porcelain)" ]; then \
#        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
#        exit 1;\
#    else \
#    	echo "\033[92mSwagger docs are in sync\033[0m";\
#    fi

###############################################################################
###                                Localnet                                 ###
###############################################################################

build-docker-go-bitsong:
	$(MAKE) -C contrib/localnet

# Run a 4-node testnet locally
localnet-start: build-linux build-docker-bitsongdnode
	@if ! [ -f build/node0/bitsongd/config/genesis.json ]; \
		then docker run --rm -v $(CURDIR)/build:/bitsongd:Z bitsongofficial/bitsongdnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 --keyring-backend=test ; \
	fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

test-docker:
	@docker build -f contrib/Dockerfile.test -t ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) .
	@docker tag ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) ${TEST_DOCKER_REPO}:$(shell git rev-parse --abbrev-ref HEAD | sed 's#/#_#g')
	@docker tag ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) ${TEST_DOCKER_REPO}:latest

test-docker-push: test-docker
	@docker push ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD)
	@docker push ${TEST_DOCKER_REPO}:$(shell git rev-parse --abbrev-ref HEAD | sed 's#/#_#g')
	@docker push ${TEST_DOCKER_REPO}:latest


########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/bitsongd -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf snapcraft-local.yaml build/

distclean: clean
	rm -rf vendor/

###############################################################################
###                                Protobuf                                 ###
###############################################################################

containerProtoVer=v0.2
containerProtoImage=tendermintdev/sdk-proto-gen:$(containerProtoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(containerProtoVer)
containerProtoGenSwagger=cosmos-sdk-proto-gen-swagger-$(containerProtoVer)
containerProtoFmt=cosmos-sdk-proto-fmt-$(containerProtoVer)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm --name $(containerProtoGen) \
		-v $(CURDIR):/workspace \
		--workdir /workspace \
		$(containerProtoImage) sh ./scripts/protocgen.sh

# This generates the SDK's custom wrapper for google.protobuf.Any. It should only be run manually when needed
proto-gen-any:
	@echo "Generating Protobuf Any"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) sh ./scripts/protocgen-any.sh

#proto-swagger-gen:
#	@echo "Generating Protobuf Swagger"
#	$(DOCKER) run --rm --name $(containerProtoGenSwagger) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) sh ./scripts/protoc-swagger-gen.sh

proto-format:
	@echo "Formatting Protobuf files"
	$(DOCKER) run --rm --name $(containerProtoFmt) \
		--user $(shell id -u):$(shell id -g) \
		-v $(CURDIR):/workspace \
		--workdir /workspace \
		tendermintdev/docker-build-proto find ./ -not -path "./third_party/*" -name *.proto -exec clang-format -i {} \;

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=master

GOGO_PROTO_URL           = https://raw.githubusercontent.com/regen-network/protobuf/cosmos
GOOGLE_PROTO_URL         = https://raw.githubusercontent.com/googleapis/googleapis/master
REGEN_COSMOS_PROTO_URL   = https://raw.githubusercontent.com/regen-network/cosmos-proto/master
COSMOS_PROTO_URL         = https://raw.githubusercontent.com/cosmos/cosmos-sdk/v0.45.4/proto/cosmos

GOGO_PROTO_TYPES         = third_party/proto/gogoproto
GOOGLE_PROTO_TYPES       = third_party/proto/google
REGEN_COSMOS_PROTO_TYPES = third_party/proto/cosmos_proto
COSMOS_PROTO_TYPES       = third_party/proto/cosmos

proto-update-deps:
	@mkdir -p $(GOGO_PROTO_TYPES)
	@curl -sSL $(GOGO_PROTO_URL)/gogoproto/gogo.proto > $(GOGO_PROTO_TYPES)/gogo.proto

	@mkdir -p $(GOOGLE_PROTO_TYPES)/api/
	@curl -sSL $(GOOGLE_PROTO_URL)/google/api/annotations.proto > $(GOOGLE_PROTO_TYPES)/api/annotations.proto
	@curl -sSL $(GOOGLE_PROTO_URL)/google/api/http.proto > $(GOOGLE_PROTO_TYPES)/api/http.proto

	@mkdir -p $(REGEN_COSMOS_PROTO_TYPES)
	@curl -sSL $(REGEN_COSMOS_PROTO_URL)/cosmos.proto > $(REGEN_COSMOS_PROTO_TYPES)/cosmos.proto

	@mkdir -p $(COSMOS_PROTO_TYPES)/base/v1beta1/
	@curl -sSL $(COSMOS_PROTO_URL)/base/v1beta1/coin.proto > $(COSMOS_PROTO_TYPES)/base/v1beta1/coin.proto

	@mkdir -p $(COSMOS_PROTO_TYPES)/base/query/v1beta1/
	@curl -sSL $(COSMOS_PROTO_URL)/base/query/v1beta1/pagination.proto > $(COSMOS_PROTO_TYPES)/base/query/v1beta1/pagination.proto

.PHONY: proto-all proto-gen proto-lint proto-check-breaking proto-update-deps

########################################
### Testing

test: test-unit
test-all: test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' -ldflags '$(ldflags)' ${PACKAGES_UNITTEST}

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

lint: golangci-lint
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./swagger/*/statik.go" -not -path "*.pb.go" | xargs gofmt -d -s
	go mod verify

benchmark:
	@go test -mod=readonly -bench=. ./...

# include simulations
# include sims.mk

.PHONY: all build-linux install install-debug \
	go-mod-cache draw-deps clean build \
	build-docker-bitsongdnode localnet-start localnet-stop test-docker test-docker-push \
	test test-all test-cover