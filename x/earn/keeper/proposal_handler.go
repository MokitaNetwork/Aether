package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mokitanetwork/aether/x/earn/types"
	aethdisttypes "github.com/mokitanetwork/aether/x/aethdist/types"
)

// HandleCommunityPoolDepositProposal is a handler for executing a passed community pool deposit proposal
func HandleCommunityPoolDepositProposal(ctx sdk.Context, k Keeper, p *types.CommunityPoolDepositProposal) error {
	fundAcc := k.accountKeeper.GetModuleAccount(ctx, aethdisttypes.FundModuleAccount)
	if err := k.distKeeper.DistributeFromFeePool(ctx, sdk.NewCoins(p.Amount), fundAcc.GetAddress()); err != nil {
		return err
	}

	err := k.DepositFromModuleAccount(ctx, aethdisttypes.FundModuleAccount, p.Amount, types.STRATEGY_TYPE_SAVINGS)
	if err != nil {
		return err
	}

	return nil

}

// HandleCommunityPoolWithdrawProposal is a handler for executing a passed community pool withdraw proposal.
func HandleCommunityPoolWithdrawProposal(ctx sdk.Context, k Keeper, p *types.CommunityPoolWithdrawProposal) error {
	// Withdraw to fund module account
	withdrawAmount, err := k.WithdrawFromModuleAccount(ctx, aethdisttypes.FundModuleAccount, p.Amount, types.STRATEGY_TYPE_SAVINGS)
	if err != nil {
		return err
	}

	// Move funds to the community pool manually
	err = k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		aethdisttypes.FundModuleAccount,
		k.distKeeper.GetDistributionAccount(ctx).GetName(),
		sdk.NewCoins(withdrawAmount),
	)
	if err != nil {
		return err
	}
	feePool := k.distKeeper.GetFeePool(ctx)
	newCommunityPool := feePool.CommunityPool.Add(sdk.NewDecCoinFromCoin(withdrawAmount))
	feePool.CommunityPool = newCommunityPool
	k.distKeeper.SetFeePool(ctx, feePool)
	return nil
}
