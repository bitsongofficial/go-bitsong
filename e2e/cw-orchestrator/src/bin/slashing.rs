extern crate abstract_client;
extern crate abstract_interface;
extern crate abstract_std;
extern crate anyhow;
extern crate bitsong_e2e_cw_orchestrator;
extern crate cosmwasm_std;
extern crate cw_orch;
extern crate cw_orch_interchain;
extern crate cw_orch_proto;

use anyhow::Result as AnyResult;

// TODO:
pub fn test_slashing_sanity() -> AnyResult<()> {
    // setup multiple nodes and validators

    //  setup delegation structure

    // have delegator slash, confirm reward is accurate

    Ok(())
}

pub fn main() {
    test_slashing_sanity().unwrap();
}
