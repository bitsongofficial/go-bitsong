###############################################################################
###                            Build & Install                              ###
###############################################################################

build-help:
	@echo "build subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make build-[command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  build-check-version                    Check Go version"
	@echo "  build-linux                            Build for Linux"
	@echo "  build-windows                          Build for Windows"
	@echo "  build-reproducible                     Build reproducible binaries"
	@echo "  build-reproducible-amd64               Build reproducible amd64 binary"
	@echo "  build-reproducible-arm64               Build reproducible arm64 binary"

# Cross-building for arm64 from amd64 (or vice-versa) takes
# a lot of time due to QEMU virtualization but it's the only way (afaik)
# to get a statically linked binary with CosmWasm
build-reproducible: build-reproducible-amd64 build-reproducible-arm64

build-check-version:
	@echo "Go version: $(GO_MAJOR_VERSION).$(GO_MINOR_VERSION)"
	@if [ $(GO_MAJOR_VERSION) -gt $(GO_MINIMUM_MAJOR_VERSION) ]; then \
		echo "Go version is sufficient"; \
		exit 0; \
	elif [ $(GO_MAJOR_VERSION) -lt $(GO_MINIMUM_MAJOR_VERSION) ]; then \
		echo '$(GO_VERSION_ERR_MSG)'; \
		exit 1; \
	elif [ $(GO_MINOR_VERSION) -lt $(GO_MINIMUM_MINOR_VERSION) ]; then \
		echo '$(GO_VERSION_ERR_MSG)'; \
		exit 1; \
	fi

build-reproducible-amd64: go.sum
	mkdir -p $(BUILDDIR)
	$(DOCKER) buildx create --name bitsongbuilder || true
	$(DOCKER) buildx use bitsongbuilder
	$(DOCKER) buildx build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
		--build-arg RUNNER_IMAGE=alpine:3.17 \
		--platform linux/amd64 \
		-t bitsong:local-amd64 \
		--load \
		-f Dockerfile .
	$(DOCKER) rm -f bitsongbinary || true
	$(DOCKER) create -ti --name bitsongbinary bitsong:local-amd64
	$(DOCKER) cp bitsongbinary:/usr/bin/bitsongd $(BUILDDIR)/bitsongd-linux-amd64
	$(DOCKER) rm -f bitsongbinary

build-reproducible-arm64: go.sum
	mkdir -p $(BUILDDIR)
	$(DOCKER) buildx create --name bitsongbuilder || true
	$(DOCKER) buildx use bitsongbuilder
	$(DOCKER) buildx build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
		--build-arg RUNNER_IMAGE=alpine:3.17 \
		--platform linux/arm64 \
		-t bitsong:local-arm64 \
		--load \
		-f Dockerfile .
	$(DOCKER) rm -f bitsongbinary || true
	$(DOCKER) create -ti --name bitsongbinary terp-core:local-arm64
	$(DOCKER) cp bitsongbinary:/usr/bin/bitsongd $(BUILDDIR)/bitsongd-linux-arm64
	$(DOCKER) rm -f bitsongbinary