#!/usr/bin/make -f

########################################
### Simulations

BINDIR ?= $(GOPATH)/bin
SIMAPP = ./app

test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.bitsongd/config/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation -Genesis=${HOME}/.bitsongd/config/genesis.json \
		-Enabled=true -NumBlocks=5 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

test-sim-fullapp:
	@echo "Running app simulation test..."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

test-sim-import-export:
	@echo "Running application import/export simulation..."
	go test -mod=readonly $(SIMAPP) -run=TestAppImportExport -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 10m