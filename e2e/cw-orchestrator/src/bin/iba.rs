// This script is used for testing a connection between 4 chains
// This script checks ibc-hook memo implementation on ibc-client
extern crate abstract_client;
extern crate abstract_interface;
extern crate abstract_std;
extern crate anyhow;
extern crate bitsong_e2e_cw_orchestrator;
extern crate cosmwasm_std;
extern crate cw_orch;
extern crate cw_orch_interchain;
extern crate cw_orch_proto;
use std::time::{SystemTime, UNIX_EPOCH};

use abstract_interface::{AccountDetails, AccountI};
use abstract_std::{
    account,
    ans_host::ExecuteMsgFns,
    ibc_client,
    ibc_host::HostAction,
    objects::{TruncatedChainId, UncheckedChannelEntry},
    IBC_CLIENT, ICS20,
};
use anyhow::Result as AnyResult;
use bitsong_e2e_cw_orchestrator::{
    bitsong_starship_interfaces, set_starship_env, BITSONG, BITSONG2,
};
use cosmwasm_std::BankMsg;
use cw_orch::prelude::*;
use cw_orch_interchain::prelude::*;
use cw_orch_proto::tokenfactory::{create_denom, get_denom, mint};

pub fn test_ibc_hook_callback() -> AnyResult<()> {
    dotenv::dotenv().ok();
    set_starship_env();
    env_logger::init();

    let starship = Starship::new(None).unwrap();
    let interchain = starship.interchain_env();

    let bitsong = interchain.get_chain(BITSONG).unwrap();

    // Create a channel between the 2 chains for the transfer ports
    // BITSONG>BITSONG2
    let bitsong_bitsong2_channel = interchain
        .create_channel(
            BITSONG,
            BITSONG2,
            &PortId::transfer(),
            &PortId::transfer(),
            "ics20-1",
            Some(cosmwasm_std::IbcOrder::Unordered),
        )?
        .interchain_channel;

    let (abstr_bitsong, _abstr_bitsong2) = bitsong_starship_interfaces(&interchain)?;

    let bitsong_sender = bitsong.sender_addr().to_string();

    // Register this channel with the abstract ibc implementation for sending tokens
    abstr_bitsong.ans_host.update_channels(
        vec![(
            UncheckedChannelEntry {
                connected_chain: TruncatedChainId::from_chain_id(BITSONG2).to_string(),
                protocol: ICS20.to_string(),
            },
            bitsong_bitsong2_channel
                .get_chain(BITSONG)?
                .channel
                .unwrap()
                .to_string(),
        )],
        vec![],
    )?;

    // Create a test account + Remote account
    let origin_account = AccountI::create_default_account(
        &abstr_bitsong,
        abstract_client::GovernanceDetails::Monarchy {
            monarch: bitsong_sender.clone(),
        },
    )?;
    origin_account.set_ibc_status(true)?;
    origin_account.create_remote_account(
        AccountDetails::default(),
        TruncatedChainId::from_chain_id(BITSONG2),
    )?;
    let test_amount: u128 = 100_000_000_000;
    let token_subdenom = format!(
        "testtoken{}",
        SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs()
    );
    // Create Denom
    create_denom(&bitsong, token_subdenom.as_str())?;
    let denom = get_denom(&bitsong, token_subdenom.as_str());
    // Mint Denom
    mint(
        &bitsong,
        &origin_account.addr_str()?,
        token_subdenom.as_str(),
        test_amount,
    )?;

    let tx_response = origin_account.execute_on_module(
        IBC_CLIENT,
        ibc_client::ExecuteMsg::RemoteAction {
            host_chain: TruncatedChainId::from_chain_id(BITSONG2),
            action: HostAction::Dispatch {
                account_msgs: vec![account::ExecuteMsg::<Empty>::Execute {
                    msgs: vec![BankMsg::Send {
                        to_address: bitsong_sender,
                        amount: vec![Coin::new(5_000_000_u128, denom.clone())],
                    }
                    .into()],
                }],
            },
        },
        vec![Coin::new(100_000_000_u128, denom.clone())],
    )?;
    interchain.await_and_check_packets(BITSONG, tx_response)?;

    let balance = bitsong.balance(&bitsong.sender_addr(), Some(denom.clone()))?;
    assert_eq!(balance, vec![Coin::new(5_000_000_u128, denom)]);
    println!("We got a callback! Result: {balance:?}");
    Ok(())
}

pub fn main() {
    test_ibc_hook_callback().unwrap();
}
