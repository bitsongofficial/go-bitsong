use abstract_interface::{Abstract, AccountI};
use abstract_std::objects::gov_type::GovernanceDetails;

use clap::Parser;
use cw_orch::daemon::networks::parse_network;
use cw_orch::prelude::*;
use interchain_bitsong_accounts::DEPLOYMENT_DAO;

pub const ABSTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[derive(Parser, Default, Debug)]
#[command(author, version, about, long_about = None)]
struct Arguments {
    /// Network Id to deploy on
    #[arg(short, long, default_value = "bitsong-2b")]
    network: String,
    /// AuthZ granter address (optional - if not provided, direct signing is used)
    #[arg(short, long, default_value = DEPLOYMENT_DAO)]
    authz_granter: Option<String>,
}

fn main() {
    env_logger::init();
    let args = Arguments::parse();
    let network = parse_network(&args.network).unwrap();

    if let Err(ref err) = init_contracts(network.into()) {
        log::error!("{}", err);
        err.chain()
            .skip(1)
            .for_each(|cause| log::error!("because: {}", cause));

        ::std::process::exit(1);
    }
}

fn init_contracts(network: ChainInfoOwned) -> anyhow::Result<()> {
    let mut chain = DaemonBuilder::new(network.clone()).build()?;
    let abs = Abstract::new(chain.clone());

    Ok(())
}
