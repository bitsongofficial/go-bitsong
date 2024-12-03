###############################################################################
###                             e2e interchain test                         ###
###############################################################################

e2e-help:
	@echo "e2e subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make [command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  e2e-basic      Run single node "
	@echo "  e2e-upgrade    Run basic planned upgrade test"
	@echo "  e2e-pfm        Run packet-forward-middleware test "

e2e: e2e-help

e2e-basic: rm-testcache
	cd e2e && go test -race -v -run TestBasicBtsgStart .

e2e-upgrade: rm-testcache
	cd e2e && go test -race -v -run TestBasicBitsongUpgrade .

e2e-pfm: rm-testcache
	cd e2e && go test -race -v -run TestPacketForwardMiddlewareRouter .

rm-testcache:
	go clean -testcache


.PHONY: test-mutation ie2e-upgrade