use abstract_interface::{Abstract, AbstractDaemonState, AccountDetails, AccountI};
use abstract_std::{native_addrs, objects::gov_type::GovernanceDetails, ACCOUNT};
use cosmwasm_std::{instantiate2_address, Binary, CanonicalAddr, Instantiate2AddressError};
use cw_blob::interface::{CwBlob, DeterministicInstantiation};
use cw_orch_daemon::{networks::BITSONG_2B, RUNTIME};

use clap::Parser;
use cw_orch::prelude::*;
use interchain_bitsong_accounts::{assert_wallet_balance, BITSONG_LOCAL_1, BITSONG_LOCAL_2};

pub const ABSTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

const CW_BLOB: &str = "cw:blob";

// Run "cargo run --example download_wasms" in the `abstract-interfaces` package before deploying!
fn init_contracts(mut networks: Vec<ChainInfoOwned>, authz_granter: &Option<String>) -> anyhow::Result<()> {
    // let networks = RUNTIME.block_on(assert_wallet_balance(networks));

    for network in networks {
        let mut chain = DaemonBuilder::new(network.clone()).build()?;
        
        // Determine admin address based on authz usage
        let admin = match authz_granter {
            Some(granter) => {
                println!("Using AuthZ granter: {}", granter);
                chain
                    .sender_mut()
                    .set_authz_granter(&Addr::unchecked(granter));
                Addr::unchecked(granter.to_string())
            },
            None => {
                println!("Using direct sender (no AuthZ)");
                chain.sender_addr()
            }
        };

        let monarch = chain.sender_addr();
        let mut abstr = Abstract::store_on(chain.clone())?;
        let mut account = AccountI::new(ACCOUNT, chain.clone());
        // set code-ids

        // account.set_default_code_id(96);
        // abstr.ans_host.set_default_code_id(93);
        // abstr.ibc.client.set_default_code_id(97);
        // abstr.ibc.host.set_default_code_id(98);
        // abstr.module_factory.set_default_code_id(95);
        // abstr.registry.set_default_code_id(94);
        // blob.set_default_code_id(92);
        // let blob_code_id = blob.code_id()?;

        abstr.ans_host.instantiate(
            &abstract_std::ans_host::InstantiateMsg {
                admin: admin.to_string(),
            },
            Some(&admin),
            &[],
        )?;

        abstr.registry.instantiate(
            &abstract_std::registry::InstantiateMsg {
                admin: admin.to_string(),
                security_enabled: Some(true),
                namespace_registration_fee: None,
            },
            Some(&admin),
            &[],
        )?;

        abstr.module_factory.instantiate(
            &abstract_std::module_factory::InstantiateMsg {
                admin: admin.to_string(),
            },
            Some(&admin),
            &[],
        )?;

        // We also instantiate ibc contracts
        abstr.ibc.instantiate(&Addr::unchecked(admin.clone()))?;

        abstr.registry.register_base(&account)?;
        abstr.registry.register_natives(abstr.contracts())?;
        abstr.registry.approve_any_abstract_modules()?;

        // Create the Abstract Account because it's needed for the fees for the dex module
        AccountI::create(
            &abstr,
            AccountDetails {
                name: "Abstract Account".to_string(),
                description: None,
                link: None,
                namespace: None,
                install_modules: vec![],
                account_id: Some(0),
            },
            GovernanceDetails::Monarchy {
                monarch: monarch.to_string(),
            },
            &[],
        )?;
    }

    // fs::copy(Path::new("~/.cw-orchestrator/state.json"), to)
    Ok(())
}

#[derive(Parser, Default, Debug)]
#[command(author, version, about, long_about = None)]
struct Arguments {
    /// AuthZ granter address (optional - if not provided, direct signing is used)
    #[arg(short, long)]
    authz_granter: Option<String>,
}

fn main() {
    // Load environment variables
    dotenv::from_path(".env").ok();
    let mnemonic = dotenv::var("LOCAL_MNEMONIC").unwrap();
    env_logger::init();
    println!("{:#?}", mnemonic);
    let args = Arguments::parse();

    let networks = vec![BITSONG_LOCAL_1.into()];

    // Determine authz usage from environment variable or command line
    let use_authz = dotenv::var("USE_AUTHZ").unwrap_or_else(|_| "false".to_string());
    let authz_granter = if use_authz == "true" && args.authz_granter.is_none() {
        // If USE_AUTHZ=true but no CLI arg provided, this is an error for init_contracts
        eprintln!("Error: USE_AUTHZ=true but --authz-granter not provided");
        std::process::exit(1);
    } else {
        args.authz_granter
    };
    
    // let dao = "bitsong13hmdq0slwmff7sej79kfa8mgnx4rl46nj2fvmlgu6u32tz6vfqesdfq4vm";

    if let Err(ref err) = init_contracts(networks, &authz_granter) {
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
