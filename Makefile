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

all: clean go-mod-cache install lint test

########################################
### CI

#ci: get_tools install lint test

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
### Testing

test: test_unit

test_unit:
	@VERSION=$(VERSION) go test $(PACKAGES_NOSIMULATION)

test_race:
	@VERSION=$(VERSION) go test -race $(PACKAGES_NOSIMULATION)

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/terra-project/core

benchmark:
	@go test -bench=. $(PACKAGES_NOSIMULATION)

lint: get_tools 


# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install clean distclean \
get_tools update_tools \
test test_cli test_unit benchmark \
build-linux  \
format update_dev_tools lint \
go-mod-cache go-sum