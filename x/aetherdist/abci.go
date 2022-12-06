package aethdist

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mokitanetwork/aether/x/aethdist/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	err := k.MintPeriodInflation(ctx)
	if err != nil {
		panic(err)
	}
}
