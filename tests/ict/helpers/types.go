package helpers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EntryPoint for various queries used throughout interchaintest library
type QueryMsg struct {
	// Tokenfactory Core
	GetConfig      *struct{}            `json:"get_config,omitempty"`
	GetBalance     *GetBalanceQuery     `json:"get_balance,omitempty"`
	GetAllBalances *GetAllBalancesQuery `json:"get_all_balances,omitempty"`

	// Unity Contract
	GetWithdrawalReadyTime *struct{} `json:"get_withdrawal_ready_time,omitempty"`

	// IBCHooks
	GetCount      *GetCountQuery      `json:"get_count,omitempty"`
	GetTotalFunds *GetTotalFundsQuery `json:"get_total_funds,omitempty"`
}

type GetBalanceQuery struct {
	// {"get_balance":{"address":"terp1...","denom":"factory/terp1.../RcqfWz"}}
	Address string `json:"address"`
	Denom   string `json:"denom"`
}

type GetAllBalancesQuery struct {
	Address string `json:"address"`
}
type GetAllBalancesResponse struct {
	// or is it wasm Coin type?
	Data []sdk.Coin `json:"data"`
}

type GetCountQuery struct {
	// {"get_total_funds":{"addr":"terp1..."}}
	Addr string `json:"addr"`
}

type GetCountResponse struct {
	// {"data":{"count":0}}
	Data *GetCountObj `json:"data"`
}

type GetCountObj struct {
	Count int64 `json:"count"`
}

type GetTotalFundsQuery struct {
	// {"get_total_funds":{"addr":"terp1..."}}
	Addr string `json:"addr"`
}

type GetTotalFundsResponse struct {
	// {"data":{"total_funds":[{"denom":"ibc/04F5F501207C3626A2C14BFEF654D51C2E0B8F7CA578AB8ED272A66FE4E48097","amount":"1"}]}}
	Data *GetTotalFundsObj `json:"data"`
}
type GetTotalFundsObj struct {
	TotalFunds []WasmCoin `json:"total_funds"`
}
type WasmCoin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
