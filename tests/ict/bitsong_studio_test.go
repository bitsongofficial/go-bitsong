package e2e

import (
	"testing"

	"github.com/CosmWasm/wasmvm/v3/types"
)

// upload studio contracts
// create an sond release
// ensure mint a song
// ensure community pool recieved cut

// Instantiates the Factory Contracts
type FactoryInitMsg struct {
	Owner                  string     `json:"owner"`
	Bs721CodeID            uint64     `json:"bs721_code_id"`
	Bs721RoyaltiesCodeID   uint64     `json:"bs721_royalties_code_id"`
	Bs721LaunchpartyCodeID uint64     `json:"bs721_launchparty_code_id"`
	Bs721CurveCodeID       uint64     `json:"bs721_curve_code_id"`
	ProtocolFeeBps         uint32     `json:"protocol_fee_bps"`
	CreateNftSaleFee       types.Coin `json:"create_nft_sale_fee"`
}

// Creates the launchparty minter contract
type MsgCreateLaunchParty struct {
	Owner                  string     `json:"owner"`
	Bs721CodeID            uint64     `json:"bs721_code_id"`
	Bs721RoyaltiesCodeID   uint64     `json:"bs721_royalties_code_id"`
	Bs721LaunchpartyCodeID uint64     `json:"bs721_launchparty_code_id"`
	Bs721CurveCodeID       uint64     `json:"bs721_curve_code_id"`
	ProtocolFeeBps         uint32     `json:"protocol_fee_bps"`
	CreateNftSaleFee       types.Coin `json:"create_nft_sale_fee"`
}

type LaunchpartyInitMsg struct {
}

type Bs721InitMsg struct {
}

type FeeSplitInitMsg struct {
}

func TestBitsongOnBitsong(t *testing.T) {
	// cfg := BaseCfg
	// numVals, numNodes := 4, 0
	// chains := CreateICTestBitsongChainCustomConfig(t, numVals, numNodes, cfg)
	// chain := chains[0].(*cosmos.CosmosChain)
	// ic, ctx, _, _ := BuildInitialChain(t, chains)

	// t.Cleanup(func() {
	// 	_ = ic.Close()
	// })

	// userFunds := sdkmath.NewInt(10_000_000_000)
	// users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chain)

	// userA, userB := users[0], users[1]

	// // store bs721
	// bs721Id, err := chain.StoreContract(ctx, userA.KeyName(), "contracts/polytone_note.wasm")
	// require.NoError(t, err)
	// // store launchpary
	// launchpartyId, err := chain.StoreContract(ctx, userA.KeyName(), "contracts/polytone_note.wasm")
	// require.NoError(t, err)
	// // store royalties
	// feeSplitterId, err := chain.StoreContract(ctx, userA.KeyName(), "contracts/polytone_note.wasm")
	// require.NoError(t, err)
	// // instantiate launchparty
	// launchpartyAddr, err := chain.InstantiateContract(ctx, user, noteId, LaunchpartyInitMsg{})
	// require.NoError(t, err)

	// // query bs721 address
	// // mint 2
	// _, err = chain.ExecuteContract(s.ctx, userA.KeyName(), launchpartyAddr, string(marshalled))
	// if err != nil {
	// 	return Callback{}, err
	// }

}
