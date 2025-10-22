use abstract_interface::{Abstract, AccountI};
use abstract_std::objects::gov_type::GovernanceDetails;

use bs_accounts::*;
use cosmwasm_std::coin;

use bs_accounts::networks::BITSONG_MAINNET;
use cw_orch::daemon::networks::parse_network;
use cw_orch::prelude::*;
use interchain_bitsong_accounts::{BITSONG_LOCAL_1, BITSONG_LOCAL_2, DEPLOYMENT_DAO};

use clap::Parser;

pub const ABSTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[derive(Parser, Default, Debug)]
#[command(author, version, about, long_about = None)]
struct Arguments {
    /// Network Id to deploy on
    #[arg(short, long, value_delimiter = ' ', num_args = 1..)]
    network_ids: Vec<String>,
    /// AuthZ granter address (optional - if not provided, direct signing is used)
    #[arg(short, long, default_value = DEPLOYMENT_DAO)]
    authz_granter: Option<String>,
    #[arg(short, long, default_value_t = false)]
    deploy_abstract: bool,
}

fn main() {
    // Load environment variables
    dotenv::from_path(".env").ok();
    env_logger::init();

    // Determine authz usage from environment variable or command line
    let args = Arguments::parse();
    let authz_granter = args.authz_granter;

    // let networks = vec![BITSONG_LOCAL_1.into(), BITSONG_LOCAL_2.into()];
    let networks = args
        .network_ids
        .iter()
        .map(|n| parse_network(n).unwrap().into())
        .collect::<Vec<_>>();

    if let Err(ref err) = full_deploy(networks, authz_granter, args.deploy_abstract) {
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

fn full_deploy(
    networks: Vec<ChainInfoOwned>,
    authz_granter: Option<String>,
    deploy_abst: bool,
) -> anyhow::Result<()> {
    for network in networks {
        let mut chain = DaemonBuilder::new(network.clone()).build()?;
        let mut admin = chain.sender_addr();
        // Conditionally set authz granter based on environment or parameter
        if let Some(granter) = &authz_granter {
            let le_granter = &Addr::unchecked(granter.to_string());
            println!("Using AuthZ granter: {}", granter);
            chain.sender_mut().set_authz_granter(le_granter);
            admin = le_granter.clone();
        } else {
            println!("Using direct sender (no AuthZ)");
        }

        // #####################################################################
        // # BITSONG ACCOUNT
        // ####################################################################
        let btsg_suite = BtsgAccountSuite::deploy_on(chain.clone(), admin.clone())?;
        let btsg_account_id = "deployment-dao";
        let fee = btsg_suite.minter.params()?.base_price;
        btsg_suite
            .account
            .approve_all(btsg_suite.market.address()?, None)?;
        btsg_suite.minter.mint_and_list(btsg_account_id)?;
        btsg_suite.minter.execute(
            &bs_accounts::Bs721AccountMinterExecuteMsgTypes::MintAndList {
                account: btsg_account_id.to_string(),
            },
            &[coin(fee.u128(), "ubtsg")],
        )?;

        // #####################################################################
        // # ABSTRACT FRAMEWORK
        // ####################################################################
        if deploy_abst {
            let deployment = match Abstract::deploy_on(chain.clone(), ()) {
                Ok(deployment) => {
                    // write_deployment(&deployment_status)?;
                    deployment
                }
                Err(e) => {
                    // write_deployment(&deployment_status)?;
                    return Err(e.into());
                }
            };

            // Create the Abstract Account because it's needed for the fees for the dex module
            AccountI::create_default_account(
                &deployment,
                GovernanceDetails::NFT {
                    collection_addr: btsg_suite.account.address()?.to_string(),
                    token_id: btsg_account_id.to_string(),
                },
            )?;
        }
    }
    // fs::copy(Path::new("~/.cw-orchestrator/state.json"), to)
    Ok(())
}
