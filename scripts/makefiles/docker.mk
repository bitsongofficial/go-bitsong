###############################################################################
###                                Docker                                  ###
###############################################################################

RUNNER_BASE_IMAGE_DISTROLESS := gcr.io/distroless/static-debian11
RUNNER_BASE_IMAGE_ALPINE := alpine:3.17
RUNNER_BASE_IMAGE_NONROOT := gcr.io/distroless/static-debian11:nonroot

docker-help:
	@echo "docker subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make [command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  docker-build                Build Docker image"
	@echo "  docker-build-distroless     Build distroless Docker image"
	@echo "  docker-build-alpine         Build alpine Docker image"
	@echo "  docker-build-nonroot        Build nonroot Docker image"
docker: docker-help

docker-build:
	@DOCKER_BUILDKIT=1 docker build \
		-t bitsong:local \
		-t bitsong:local-distroless \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg RUNNER_IMAGE=$(RUNNER_BASE_IMAGE_DISTROLESS) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
		-f Dockerfile .

docker-build-distroless: docker-build

docker-build-alpine:
	@DOCKER_BUILDKIT=1 docker build \
		-t bitsong:local-alpine \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg RUNNER_IMAGE=$(RUNNER_BASE_IMAGE_ALPINE) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
		-f Dockerfile .

docker-build-nonroot:
	@DOCKER_BUILDKIT=1 docker build \
		-t bitsong:local-nonroot \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg RUNNER_IMAGE=$(RUNNER_BASE_IMAGE_NONROOT) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
		-f Dockerfile .
