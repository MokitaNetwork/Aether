package aethdist

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/mokitanetwork/aether/x/aethdist/keeper"
	"github.com/mokitanetwork/aether/x/aethdist/types"
)

// NewCommunityPoolMultiSpendProposalHandler
func NewCommunityPoolMultiSpendProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.CommunityPoolMultiSpendProposal:
			return keeper.HandleCommunityPoolMultiSpendProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized aethdist proposal content type: %T", c)
		}
	}
}
