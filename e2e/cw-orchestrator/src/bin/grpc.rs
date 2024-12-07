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
pub fn test_single_node_grpc_sanity() -> AnyResult<()> {
    // setup single localbitsong instance and query grpc
    //  perform all grpc requests from modules available, including ibc queries with grpc
    Ok(())
}

pub fn main() {
    test_single_node_grpc_sanity().unwrap();
}
