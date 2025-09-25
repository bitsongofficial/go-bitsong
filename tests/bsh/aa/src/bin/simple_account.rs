use abstract_interface::{Abstract, AccountI};
use abstract_std::objects::gov_type::GovernanceDetails;
use cw_orch_daemon::networks::BITSONG_2B;

use cw_orch::prelude::*;
use interchain_bitsong_accounts::DEPLOYMENT_DAO;

pub const ABSTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

// Run "cargo run --example download_wasms" in the `abstract-interfaces` package before deploying!
fn init_contracts(networks: Vec<ChainInfoOwned>) -> anyhow::Result<()> {
    // let networks = RUNTIME.block_on(assert_wallet_balance(networks));

    for network in networks {
        let mut chain = DaemonBuilder::new(network.clone()).build()?;

        // use Authz granted by Bitsong Deployment SubDAO
        chain
            .sender_mut()
            .set_authz_granter(&Addr::unchecked(DEPLOYMENT_DAO));

        let monarch = chain.sender_addr();
        let abstr = Abstract::load_from(chain.clone())?;

        // Create the Abstract Account because it's needed for the fees for the dex module
        AccountI::create_default_account(
            &abstr,
            GovernanceDetails::Monarchy {
                monarch: monarch.to_string(),
            },
        )?;
    }

    // fs::copy(Path::new("~/.cw-orchestrator/state.json"), to)
    Ok(())
}

fn main() {
    dotenv::from_path(".env").ok();
    env_logger::init();

    let networks = vec![BITSONG_2B.into()];
    if let Err(ref err) = init_contracts(networks) {
        log::error!("{}", err);
        err.chain()
            .skip(1)
            .for_each(|cause| log::error!("because: {}", cause));

        // The backtrace is not always generated. Try to run this example
        // with `$env:RUST_BACKTRACE=1`.
        //    if let Some(backtrace) = e.backtrace() {
        //        log::debug!("backtrace: {:?}", backtrace);
        //    }

        ::std::process::exit(1);
    }
}
