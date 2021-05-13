package types

import (
	fmt "fmt"
	"math/big"

	"github.com/gogo/protobuf/proto"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const FanTokenDecimal = 6

var (
	_      proto.Message = &FanToken{}
	tenInt               = big.NewInt(10)
)

// TokenI defines an interface for Token
type FanTokenI interface {
	GetSymbol() string
	GetDenom() string
	GetName() string
	GetMaxSupply() sdk.Int
	GetMintable() bool
	GetOwner() sdk.AccAddress

	ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error)
	ToMinCoin(coin sdk.DecCoin) (sdk.Coin, error)
}

// NewToken constructs a new Token instance
func NewFanToken(
	symbol string,
	name string,
	maxSupply sdk.Int,
	mintable bool,
	descriptioin string,
	owner sdk.AccAddress,
) FanToken {
	return FanToken{
		Symbol:      symbol,
		Name:        name,
		MaxSupply:   maxSupply,
		Mintable:    mintable,
		Description: descriptioin,
		Owner:       owner.String(),
	}
}

// GetSymbol implements exported.TokenI
func (t FanToken) GetSymbol() string {
	return t.Symbol
}

// GetDenom implements exported.TokenI
func (t FanToken) GetDenom() string {
	return fmt.Sprintf("%s%s", "u", t.Symbol)
}

// GetName implements exported.TokenI
func (t FanToken) GetName() string {
	return t.Name
}

// GetMaxSupply implements exported.TokenI
func (t FanToken) GetMaxSupply() sdk.Int {
	return t.MaxSupply
}

// GetMintable implements exported.TokenI
func (t FanToken) GetMintable() bool {
	return t.Mintable
}

// GetOwner implements exported.TokenI
func (t FanToken) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(t.Owner)
	return owner
}

func (t FanToken) String() string {
	bz, _ := yaml.Marshal(t)
	return string(bz)
}

// ToMainCoin returns the main denom coin from args
func (t FanToken) ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error) {
	if t.Symbol != coin.Denom && t.GetDenom() != coin.Denom {
		return sdk.NewDecCoinFromDec(coin.Denom, sdk.ZeroDec()), sdkerrors.Wrapf(ErrTokenNotExists, "token not match")
	}

	if t.Symbol == coin.Denom {
		return sdk.NewDecCoin(coin.Denom, coin.Amount), nil
	}

	precision := new(big.Int).Exp(tenInt, big.NewInt(FanTokenDecimal), nil)
	// dest amount = src amount / 10^(scale)
	amount := sdk.NewDecFromInt(coin.Amount).Quo(sdk.NewDecFromBigInt(precision))
	return sdk.NewDecCoinFromDec(t.Symbol, amount), nil
}

// ToMinCoin returns the min denom coin from args
func (t FanToken) ToMinCoin(coin sdk.DecCoin) (newCoin sdk.Coin, err error) {
	if t.Symbol != coin.Denom && t.GetDenom() != coin.Denom {
		return sdk.NewCoin(coin.Denom, sdk.ZeroInt()), sdkerrors.Wrapf(ErrTokenNotExists, "token not match")
	}

	if t.GetDenom() == coin.Denom {
		return sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()), nil
	}

	precision := new(big.Int).Exp(tenInt, big.NewInt(FanTokenDecimal), nil)
	// dest amount = src amount * 10^(dest scale)
	amount := coin.Amount.Mul(sdk.NewDecFromBigInt(precision))
	return sdk.NewCoin(t.GetDenom(), amount.TruncateInt()), nil
}
