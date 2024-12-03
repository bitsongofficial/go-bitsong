package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	e "github.com/pkg/errors"
)

const (
	CoinExponent     = 6
	CoinType         = 639
	CoinUnit         = "btsg"
	MicroCoinUnit    = "ubtsg"
	DefaultBondDenom = MicroCoinUnit
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = "bitsong"
)

var (
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32PrefixAccAddr + "pub"
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32PrefixAccAddr + "valoper"
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32PrefixAccAddr + "valoperpub"
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32PrefixAccAddr + "valcons"
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32PrefixAccAddr + "valconspub"
)

func init() {
	SetAddressPrefixes()
	RegisterTokenDenomination()
}

func RegisterTokenDenomination() {
	err := sdk.RegisterDenom(CoinUnit, sdk.OneDec())
	if err != nil {
		panic(err)
	}
	err = sdk.RegisterDenom(MicroCoinUnit, sdk.NewDecWithPrec(1, CoinExponent))
	if err != nil {
		panic(err)
	}
}

func SetAddressPrefixes() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	// 639 is the registered coin type for BTSG
	// Following the coin type registered at https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	config.SetCoinType(CoinType)
	config.SetPurpose(CoinExponent)

	// This is copied from the cosmos sdk v0.43.0-beta1
	// source: https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/types/address.go#L141
	config.SetAddressVerifier(func(bytes []byte) error {
		if len(bytes) == 0 {
			return e.Wrap(sdkerrors.ErrUnknownAddress, "addresses cannot be empty")
		}

		if len(bytes) > address.MaxAddrLen {
			e.Wrapf(sdkerrors.ErrUnknownAddress, "address max length is %d, got %d", address.MaxAddrLen, len(bytes))
		}

		// TODO: Do we want to allow addresses of lengths other than 20 and 32 bytes?
		if len(bytes) != 20 && len(bytes) != 32 {
			return e.Wrapf(sdkerrors.ErrUnknownAddress, "address length must be 20 or 32 bytes, got %d", len(bytes))
		}

		return nil
	})
}
