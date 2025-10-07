###############################################################################
###                                  Proto                                  ###
###############################################################################

protoVer=0.17.1
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)
SWAGGER_DIR=./swagger-proto
THIRD_PARTY_DIR=$(SWAGGER_DIR)/third_party

proto-help:
	@echo "proto subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make proto-[command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  proto-all              Run proto-format and proto-gen"
	@echo "  proto-check-breaking   Check breaking instances"
	@echo "  proto-gen-swagger     Generate Protobuf files"
	@echo "  proto-gen-pulsar       Generate Protobuf files"
	@echo "  proto-format           Format Protobuf files"
	@echo "  proto-lint             Lint Protobuf files"
	@echo "  proto-image-build      Build the protobuf Docker image"
	@echo "  proto-image-push       Push the protobuf Docker image"
	@echo "  proto-docs             Create Swagger API docs"

proto: proto-help

proto-all: proto-format proto-gen-swagger

proto-gen-swagger:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/gen-swagger.sh

proto-gen-pulsar:
	@echo "Generating Dep-Inj Protobuf files"
	@$(protoImage) sh ./scripts/gen-pulsar.sh

# linux only
proto-format:
	@echo "Formatting Protobuf files"
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./proto -name "*.proto" -exec clang-format -i {} \;

proto-update-deps:
	@echo "Updating Protobuf dependencies"
	@$(DOCKER) run --rm -v $(CURDIR)/proto:/workspace --workdir /workspace $(protoImageName) buf mod update

proto-lint:
	@$(DOCKER) run --rm -v $(CURDIR)/proto:/workspace --workdir /workspace $(protoImageName) buf lint --error-format=json

proto-check-breaking:
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main
