###############################################################################
###                                  Proto                                  ###
###############################################################################

proto-help:
	@echo "proto subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make proto-[command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  proto-all        Run proto-format and proto-gen"
	@echo "  proto-gen        Generate Protobuf files"
	@echo "  proto-format     Format Protobuf files"
	@echo "  proto-image-build  Build the protobuf Docker image"
	@echo "  proto-image-push  Push the protobuf Docker image"
	@echo "  proto-docs        Create Swagger API docs"

proto: proto-help
proto-all: proto-format proto-gen

PROTO_BUILDER_IMAGE=ghcr.io/cosmos/proto-builder:0.14.0
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(PROTO_BUILDER_IMAGE)

proto-all: proto-format proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(DOCKER) run --rm -u 0 -v $(CURDIR):/workspace --workdir /workspace $(PROTO_BUILDER_IMAGE) sh ./scripts/protocgen.sh

# linux only
proto-format:
	@echo "Formatting Protobuf files"
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./proto -name "*.proto" -exec clang-format -i {} \;



SWAGGER_DIR=./swagger-proto
THIRD_PARTY_DIR=$(SWAGGER_DIR)/third_party

proto-download-deps:
	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-sdk.git" && \
	git config core.sparseCheckout true && \
	printf "proto\nthird_party\n" > .git/info/sparse-checkout && \
	git pull origin main && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/ibc-go.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin main && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ibc_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-proto.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin main && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_proto_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/gogoproto" && \
	curl -SSL https://raw.githubusercontent.com/cosmos/gogoproto/main/gogoproto/gogo.proto > "$(THIRD_PARTY_DIR)/gogoproto/gogo.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/google/api" && \
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > "$(THIRD_PARTY_DIR)/google/api/annotations.proto"
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > "$(THIRD_PARTY_DIR)/google/api/http.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos/ics23/v1" && \
	curl -sSL https://raw.githubusercontent.com/cosmos/ics23/master/proto/cosmos/ics23/v1/proofs.proto > "$(THIRD_PARTY_DIR)/cosmos/ics23/v1/proofs.proto"



proto-docs:
	@echo
	@echo "=========== Generate Message ============"
	@echo
	sh ./scripts/protoc-swagger-gen.sh

	@echo
	@echo "=========== Generate Complete ============"
	@echo
.PHONY: docs