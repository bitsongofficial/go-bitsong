###############################################################################
###                                Localnet                                 ###
###############################################################################
#
# Please refer to https://github.com/bitsongofficial/go-bitsong/blob/main/ict/localbitsong/README.md for detailed
# usage of localnet.

localnet-help:
	@echo "build subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make [command]"
	@echo ""
	@echo "Available Commands:"
	@echo "localnet-start  		 Start localnet"
	@echo "localnet-stop  		 Stop localnet"
	@echo "test-docker-push  	Push testnet docker image"

localnet: localnet-help

localnet-keys:
	. ict/localbitsong/scripts/add_keys.sh

localnet-init: localnet-clean localnet-build

localnet-build:
	@DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker compose -f ict/localbitsong/docker-compose.yml build

localnet-start:
	@STATE="" docker compose -f ict/localbitsong/docker-compose.yml up

localnet-startd:
	@STATE="" docker compose -f ict/localbitsong/docker-compose.yml up -d

localnet-stop:
	@STATE="" docker compose -f ict/localbitsong/docker-compose.yml down

localnet-clean:
	@rm -rfI $(HOME)/.bitsongd-local/