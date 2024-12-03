###############################################################################
###                                Localnet                                 ###
###############################################################################

localnet-help:
	@echo "build subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make [command]"
	@echo ""
	@echo "Available Commands:"
	@echo "localnet-start  		 Start a 4-node testnet locally"
	@echo "localnet-stop  		 Stop  a local testnet"
	@echo "test-docker-push  	Push testnet docker image"

localnet: localnet-help

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