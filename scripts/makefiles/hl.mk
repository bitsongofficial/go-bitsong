###############################################################################
###             heighliner - (used for docker containers)                   ###
###############################################################################
.PHONY: heighliner-get heighliner-local-image heighliner

heighliner-help:
	@echo "heighliner subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make heighliner-[command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  hl-get        	 Install Heighliner"
	@echo "  hl-local-image    Create a local image"
	@echo "  hl-previous-image    Create a local image from the previous version from upgrade."
	@echo ""
	@echo ""


hl: heighliner-help

hl-get:
	git clone https://github.com/strangelove-ventures/heighliner.git
	cd heighliner && go install

hl-local-image:
ifeq (,$(shell which heighliner))
	echo 'heighliner' binary not found. Consider running `make hl-get`
else 
	heighliner build -c bitsong  -o bitsongofficial/go-bitsong --local -f ./chains.yaml -t v0.18.0
endif