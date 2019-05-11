VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
GOTOOLS = \
	github.com/rakyll/statik
GOBIN ?= $(GOPATH)/bin
SHASUM := $(shell which sha256sum)

export GO111MODULE = on

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
    GCC = $(shell command -v gcc 2> /dev/null)
    ifeq ($(GCC),)
      $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif

# process linker flags

ldflags = -X github.com/BitSongOfficial/go-bitsong/version.Version=$(VERSION) \
					-X github.com/BitSongOfficial/go-bitsong/version.Commit=$(COMMIT) \
					-X "github.com/BitSongOfficial/go-bitsong/version.BuildTags=$(build_tags)" \

ifneq ($(SHASUM),)
	ldflags += -X github.com/BitSongOfficial/go-bitsong/version.GoSumHash=$(shell sha256sum go.sum | cut -d ' ' -f1)
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

########################################
### All

all: clean go-mod-cache install lint

########################################
### Build/Install

build: 
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/bitsongd.exe ./cmd/bitsongd
	go build $(BUILD_FLAGS) -o build/bitsongcli.exe ./cmd/bitsongcli
else
	go build $(BUILD_FLAGS) -o build/bitsongd ./cmd/bitsongd
	go build $(BUILD_FLAGS) -o build/bitsongcli ./cmd/bitsongcli
endif

build-linux:
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

install:
	go install $(BUILD_FLAGS) ./cmd/bitsongd
	go install $(BUILD_FLAGS) ./cmd/bitsongcli

########################################
### Tools & dependencies

get_tools:
	go get github.com/rakyll/statik
	go get github.com/golangci/golangci-lint/cmd/golangci-lint

update_tools:
	@echo "--> Updating tools to correct version"
	$(MAKE) --always-make get_tools

go-mod-cache: go-sum
	@echo "--> Download go modules to local cache"
	@go mod download

go-sum: get_tools
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

go-release:
	@echo "--> Dry run for go-release"
	BUILD_TAGS=$(shell echo \"$(build_tags)\") GOSUM=$(shell sha256sum go.sum | cut -d ' ' -f1) goreleaser release --skip-publish --rm-dist --debug

clean:
	rm -rf ./dist
	rm -rf ./build

distclean: clean
	rm -rf vendor/


########################################
### Local validator nodes using docker and docker-compose

build-docker-terradnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/terrad/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/terrad:Z tendermint/terradnode testnet --v 5 -o . --starting-ip-address 192.168.10.2; fi
	# replace docker ip to local port, mapped
	sed -i -e 's/192.168.10.2:26656/localhost:26656/g; s/192.168.10.3:26656/localhost:26659/g; s/192.168.10.4:26656/localhost:26661/g; s/192.168.10.5:26656/localhost:26663/g' $(CURDIR)/build/node4/terrad/config/config.toml
	# change allow duplicated ip option to prevent the error : cant not route ~
	sed -i -e 's/allow_duplicate_ip \= false/allow_duplicate_ip \= true/g' `find $(CURDIR)/build -name "config.toml"`
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down


# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install clean distclean \
get_tools update_tools \
build-linux  \
format update_dev_tools lint \
go-mod-cache go-sum