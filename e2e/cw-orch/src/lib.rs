extern crate abstract_client;
extern crate abstract_interface;
extern crate abstract_std;
extern crate anyhow;
extern crate cosmwasm_std;
extern crate cw_orch;
extern crate cw_orch_interchain;

use abstract_interface::{
    connection::connect_one_way_to, Abstract, AccountDetails, AccountI, AccountQueryFns,
};
use abstract_std::objects::{AccountId, AccountTrace, TruncatedChainId};
use anyhow::Result as AnyResult;
use cosmwasm_std::{coins, Coin};
use cw_orch::prelude::*;
use cw_orch_interchain::{
    core::{IbcQueryHandler, InterchainEnv},
    daemon::DaemonInterchain,
    prelude::Starship,
};
use networks::{ChainKind, NetworkInfo};

pub const BITSONG: &str = "bitsong-1";
pub const STARGAZE: &str = "stargaze-1";
pub const OSMOSIS: &str = "osmosis-1";

// Note: Truncated chain id have to be different
pub const BITSONG2: &str = "bitsongtwo-1";

pub const TEST_ACCOUNT_NAME: &str = "account-test";
pub const TEST_ACCOUNT_DESCRIPTION: &str = "Description of an account";
pub const TEST_ACCOUNT_LINK: &str = "https://skeret.jeret";

pub fn set_env() {
    std::env::set_var("STATE_FILE", "daemon_state.json"); // Set in code for tests
    std::env::set_var("ARTIFACTS_DIR", "../artifacts"); // Set in code for tests
}
// ANCHOR: osmosis
pub const BITSONG_NETWORK: NetworkInfo = NetworkInfo {
    chain_name: "bitsong",
    pub_address_prefix: "bitsong",
    coin_type: 118u32,
};

pub const LOCALBITSONG: ChainInfo = ChainInfo {
    kind: ChainKind::Local,
    chain_id: "local-1",
    gas_denom: "ubtsg",
    gas_price: 0.025,
    grpc_urls: &["https://localhost:8080"],
    network_info: BITSONG_NETWORK,
    lcd_url: None,
    fcd_url: None,
};

// STARSHIP

// Set in code for starship tests
pub fn set_starship_env() {
    std::env::set_var("STATE_FILE", "starship-state.json");
    std::env::set_var("ARTIFACTS_DIR", "../artifacts");
}

pub fn create_test_remote_account<Chain: IbcQueryHandler, IBC: InterchainEnv<Chain>>(
    abstr_origin: &Abstract<Chain>,
    origin_id: &str,
    remote_id: &str,
    interchain: &IBC,
    funds: Vec<Coin>,
) -> anyhow::Result<(AccountI<Chain>, AccountId)> {
    let origin_name = TruncatedChainId::from_chain_id(origin_id);
    let remote_name = TruncatedChainId::from_chain_id(remote_id);

    // Create a local account for testing
    let account_name = TEST_ACCOUNT_NAME.to_string();
    let description = Some(TEST_ACCOUNT_DESCRIPTION.to_string());
    let link = Some(TEST_ACCOUNT_LINK.to_string());
    let origin_account = AccountI::create(
        abstr_origin,
        AccountDetails {
            name: account_name.clone(),
            description: description.clone(),
            link: link.clone(),
            install_modules: vec![],
            namespace: None,
            account_id: None,
        },
        abstract_std::objects::gov_type::GovernanceDetails::Monarchy {
            monarch: abstr_origin
                .registry
                .environment()
                .sender_addr()
                .to_string(),
        },
        &funds,
    )?;

    // We need to enable ibc on the account.
    origin_account.set_ibc_status(true)?;

    // Now we send a message to the client saying that we want to create an account on the
    // host chain
    let register_tx = origin_account.register_remote_account(remote_name)?;

    interchain.await_and_check_packets(origin_id, register_tx)?;

    // After this is all ended, we return the account id of the account we just created on the remote chain
    let account_config = origin_account.config()?;
    let remote_account_id = AccountId::new(
        account_config.account_id.seq(),
        AccountTrace::Remote(vec![origin_name]),
    )?;

    Ok((origin_account, remote_account_id))
}

pub fn bitsong_starship_interfaces(
    interchain: &DaemonInterchain<Starship>,
) -> AnyResult<(Abstract<Daemon>, Abstract<Daemon>)> {
    let bitsong = interchain.get_chain(BITSONG).unwrap();
    let bitsong2 = interchain.get_chain(BITSONG2).unwrap();
    // Just return if already deployed
    if let Ok(bitsong_deployment) = Abstract::load_from(bitsong.clone()) {
        return Ok((bitsong_deployment, Abstract::load_from(bitsong2)?));
    }
    // Deploy and connect if not deployed yet

    // Send some funds for deploying abstract
    bitsong.rt_handle.block_on(bitsong.sender().bank_send(
        &bitsong.sender_addr(),
        coins(10_000_000_000_000, bitsong.chain_info().gas_denom.clone()),
    ))?;
    bitsong2.rt_handle.block_on(bitsong2.sender().bank_send(
        &bitsong2.sender_addr(),
        coins(10_000_000_000_000, bitsong2.chain_info().gas_denom.clone()),
    ))?;
    let abstr_bitsong = Abstract::deploy_on(bitsong.clone(), bitsong.sender().clone())?;
    let abstr_bitsong2 = Abstract::deploy_on(bitsong2.clone(), bitsong.sender().clone())?;
    connect_one_way_to(&abstr_bitsong, &abstr_bitsong2, interchain)?;

    Ok((abstr_bitsong, abstr_bitsong2))
}
