###############################################################################
###                             ict interchain test                         ###
###############################################################################

ict-help:
	@echo "ict subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make [command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  ict-basic         Test all core cosmos-sdk & bitsong modules are operational"
	@echo "  ict-pfm           Test packet-forward-middleware is operational"
	@echo "  ict-polytone      Test wasm + IBC support is operational"
	@echo "  ict-staking-hooks Test staking hooks are operational"
	@echo "  ict-ibc-hooks     Test IBC-Hooks are operational"
# 	@echo "  ict-upgrade    Run basic planned upgrade test"

ict: ict-help

ict-basic: rm-testcache
	cd tests/ict && go test -race -v -run TestBasicBtsgStart .

ict-upgrade: rm-testcache
	cd tests/ict && go test -race -v -run TestBasicBitsongUpgrade .

ict-pfm: rm-testcache
	cd tests/ict && go test -race -v -run TestPacketForwardMiddlewareRouter .

ict-polytone: rm-testcache
	cd tests/ict && go test -race -v -run TestPolytoneOnBitsong .

ict-ibc-hooks: rm-testcache
	cd tests/ict && go test -race -v -run TestBtsgIBCHooks .
	
ict-staking-hooks: rm-testcache
	cd tests/ict && go test -race -v -run TestPolytoneOnBitsong .


# ict-slashing: rm-testcache
# 	cd tests/ict && go test -race -v -run TestBasicBitsongSlashing .

rm-testcache:
	go clean -testcache


.PHONY: test-mutation ie2e-upgrade