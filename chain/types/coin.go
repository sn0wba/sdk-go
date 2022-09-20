package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// INJ defines the default coin denomination used in Ethermint in:
	//
	// - Staking parameters: denomination used as stake in the dPoS chain
	// - Mint parameters: denomination minted due to fee distribution rewards
	// - Governance parameters: denomination used for spam prevention in proposal deposits
	// - Crisis parameters: constant fee denomination used for spam prevention to check broken invariant
	// - EVM parameters: denomination used for running EVM state transitions in Ethermint.
	GotabitCoin string = "ugtb"

	// BaseDenomUnit defines the base denomination unit for Photons.
	// 1 photon = 1x10^{BaseDenomUnit} ugtb
	BaseDenomUnit = 6
)

// NewGotabitCoin is a utility function that returns an "ugtb" coin with the given sdk.Int amount.
// The function will panic if the provided amount is negative.
func NewGotabitCoin(amount sdk.Int) sdk.Coin {
	return sdk.NewCoin(GotabitCoin, amount)
}

// NewGotabitDecCoin is a utility function that returns an "ugtb" decimal coin with the given sdk.Int amount.
// The function will panic if the provided amount is negative.
func NewGotabitDecCoin(amount sdk.Int) sdk.DecCoin {
	return sdk.NewDecCoin(GotabitCoin, amount)
}

// NewGotabitCoinInt64 is a utility function that returns an "ugtb" coin with the given int64 amount.
// The function will panic if the provided amount is negative.
func NewGotabitCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(GotabitCoin, amount)
}
