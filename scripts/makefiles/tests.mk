

test-help:
	@echo "test subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make [command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  ict                	 View ict tests available to run"
	@echo "  test-all                Run all tests"
	@echo "  test-unit               Run unit tests"
	@echo "  test-benchmark          Run benchmark tests"
	@echo "  test-cover              Run coverage tests"
	@echo "  test-race               Run race tests"
	@echo "  test-race               Run race tests"
	@echo "  test-ict                View current InterchainTests configured"
	@echo "  test-bsh                View current bash script options"

test: test-help
test-all: test-race test-cover test-unit

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' -ldflags '$(ldflags)' ${PACKAGES_UNITTEST}

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

test-benchmark:
	@go test -mod=readonly -bench=. ./...

test-ict:
	@make ict

test-bsh:
	@make bsh

# include simulations
# include sims.mk