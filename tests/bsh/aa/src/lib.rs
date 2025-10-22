use cw_orch::environment::ChainKind;
use cw_orch::prelude::{networks::bitsong::BITSONG_NETWORK, *};

pub const DEPLOYMENT_DAO: &str =
    "bitsong13hmdq0slwmff7sej79kfa8mgnx4rl46nj2fvmlgu6u32tz6vfqesdfq4vm";
const GAS_TO_DEPLOY: u64 = 60_000_000;

pub const BITSONG_LOCAL_1: ChainInfo = ChainInfo {
    kind: ChainKind::Local,
    chain_id: "test-1",
    gas_denom: "ubtsg",
    gas_price: 0.075,
    grpc_urls: &["http://localhost:9090"],
    network_info: BITSONG_NETWORK,
    lcd_url: Some("http://localhost:1317"),
    fcd_url: None,
};

pub const BITSONG_LOCAL_2: ChainInfo = ChainInfo {
    kind: ChainKind::Local,
    chain_id: "test-2",
    gas_denom: "ubtsg",
    gas_price: 0.075,
    grpc_urls: &["http://localhost:10090"],
    network_info: BITSONG_NETWORK,
    lcd_url: Some("http://localhost:1318"),
    fcd_url: None,
};

pub async fn assert_wallet_balance(mut chains: Vec<ChainInfoOwned>) -> Vec<ChainInfoOwned> {
    // check that the wallet has enough gas on all the chains we want to support
    for chain_info in &chains {
        let chain = DaemonAsyncBuilder::new(chain_info.clone())
            .build()
            .await
            .unwrap();

        println!("{:#?}", chain.sender_addr());
        let gas_denom = chain.state().chain_data.gas_denom.clone();
        let gas_price = chain.state().chain_data.gas_price;
        let fee = (GAS_TO_DEPLOY as f64 * gas_price) as u128;
        let bank = queriers::Bank::new_async(chain.channel());
        let balance = bank
            ._balance(&chain.sender_addr(), Some(gas_denom.clone()))
            .await
            .unwrap()
            .clone()[0]
            .clone();

        log::debug!(
            "Checking balance {} on chain {}, address {}. Expecting {}{}",
            balance.amount,
            chain_info.chain_id,
            chain.sender_addr(),
            fee,
            gas_denom
        );
        if fee > balance.amount.u128() {
            panic!("Not enough funds on chain {} to deploy the contract. Needed: {}{} but only have: {}{}", chain_info.chain_id, fee, gas_denom, balance.amount, gas_denom);
        }
        // check if we have enough funds
    }

    chains
}
