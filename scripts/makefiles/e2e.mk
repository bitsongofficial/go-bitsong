###############################################################################
###                             e2e interchain test                         ###
###############################################################################

e2e-help:
	@echo "e2e subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make e2e-[command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  e2e-basic      Run single node "
	@echo "  e2e-upgrade    Run basic planned upgrade test"

e2e: e2e-help

e2e-basic: rm-testcache
	cd e2e && go test -race -v -run TestBasicBitsongStart .

e2e-upgrade: rm-testcache
	cd e2e && go test -race -v -run TestBasicBitsongUpgrade .


rm-testcache:
	go clean -testcache


.PHONY: test-mutation ie2e-upgrade