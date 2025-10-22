bsh-help:
	@echo "Bitsong e2e testing subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make bsh-[command]"
	@echo ""
	@echo "Testing bitsong app functions using its daemon + sh scripts"
	@echo "Available Commands:"
	@echo "  bsh                	 View bitsong sh tests available to run"
	@echo "  bsh-all 			 	 Run all sh tests in repo"
	@echo "  bsh-nfts 			 	 Run sh test for x/nft module"
	@echo "  bsh-ibchook 			 Run sh test for ibc hook sanity"
	@echo "  bsh-pfm 		     	 Run sh test for packet-forward-middleware sanity"
	@echo "  bsh-aa 		     	 Run sh test for sane deployment & use of Abstract Account"
	@echo "  bsh-polytone 			 Run sh test for ibc + wasm sanity"
	@echo "  bsh-staking-hooks 		 Run sh test for staking hook sanity"
	@echo "  bsh-upgrade 		     Run sh test for upgrade proposal & performance sanity"

bsh: bsh-help
bsh-all: bsh-upgrade bsh-staking-hooks bsh-polytone bsh-aa bsh-pfm bsh-ibchook bsh-nfts
bsh-aa: 
	cd tests/bsh/aa && sh a.sh
bsh-ibchook: 
	cd tests/bsh/ibchook && sh a.sh
bsh-upgrade: 
	cd tests/bsh/upgrade && sh a.sh
bsh-staking-hooks: 
	cd tests/bsh/staking-hooks && sh a.sh
bsh-polytone: 
	cd tests/bsh/polytone && sh a.sh
bsh-pfm: 
	cd tests/bsh/pfm && sh a.sh
bsh-nfts: 
	cd tests/bsh/nft && sh a.sh


# include simulations
# include sims.mk