package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	earnkeeper "github.com/mokitanetwork/aether/x/earn/keeper"
	liquidkeeper "github.com/mokitanetwork/aether/x/liquid/keeper"
	liquidtypes "github.com/mokitanetwork/aether/x/liquid/types"
	savingskeeper "github.com/mokitanetwork/aether/x/savings/keeper"
)

var _ govtypes.TallyHandler = TallyHandler{}

// TallyHandler is the tally handler for aeth
type TallyHandler struct {
	gk  govkeeper.Keeper
	stk stakingkeeper.Keeper
	svk savingskeeper.Keeper
	ek  earnkeeper.Keeper
	lk  liquidkeeper.Keeper
	bk  bankkeeper.Keeper
}

// NewTallyHandler creates a new tally handler.
func NewTallyHandler(
	gk govkeeper.Keeper, stk stakingkeeper.Keeper, svk savingskeeper.Keeper,
	ek earnkeeper.Keeper, lk liquidkeeper.Keeper, bk bankkeeper.Keeper,
) TallyHandler {
	return TallyHandler{
		gk:  gk,
		stk: stk,
		svk: svk,
		ek:  ek,
		lk:  lk,
		bk:  bk,
	}
}

func (th TallyHandler) Tally(ctx sdk.Context, proposal types.Proposal) (passes bool, burnDeposits bool, tallyResults types.TallyResult) {
	results := make(map[types.VoteOption]sdk.Dec)
	results[types.OptionYes] = sdk.ZeroDec()
	results[types.OptionAbstain] = sdk.ZeroDec()
	results[types.OptionNo] = sdk.ZeroDec()
	results[types.OptionNoWithVeto] = sdk.ZeroDec()

	totalVotingPower := sdk.ZeroDec()
	currValidators := make(map[string]types.ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	th.stk.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		currValidators[validator.GetOperator().String()] = types.NewValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			sdk.ZeroDec(),
			types.WeightedVoteOptions{},
		)

		return false
	})

	th.gk.IterateVotes(ctx, proposal.ProposalId, func(vote types.Vote) bool {
		// if validator, just record it in the map
		voter, err := sdk.AccAddressFromBech32(vote.Voter)

		if err != nil {
			panic(err)
		}

		valAddrStr := sdk.ValAddress(voter.Bytes()).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Options
			currValidators[valAddrStr] = val
		}

		// iterate over all delegations from voter, deduct from any delegated-to validators
		th.stk.IterateDelegations(ctx, voter, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
			valAddrStr := delegation.GetValidatorAddr().String()

			if val, ok := currValidators[valAddrStr]; ok {
				// There is no need to handle the special case that validator address equal to voter address.
				// Because voter's voting power will tally again even if there will deduct voter's voting power from validator.
				val.DelegatorDeductions = val.DelegatorDeductions.Add(delegation.GetShares())
				currValidators[valAddrStr] = val

				// delegation shares * bonded / total shares
				votingPower := delegation.GetShares().MulInt(val.BondedTokens).Quo(val.DelegatorShares)

				for _, option := range vote.Options {
					subPower := votingPower.Mul(option.Weight)
					results[option.Option] = results[option.Option].Add(subPower)
				}
				totalVotingPower = totalVotingPower.Add(votingPower)
			}

			return false
		})

		// get voter baeth and update total voting power and results
		addrBaeth := th.getAddrBaeth(ctx, voter).toCoins()
		for _, coin := range addrBaeth {
			valAddr, err := liquidtypes.ParseLiquidStakingTokenDenom(coin.Denom)
			if err != nil {
				break
			}

			// reduce delegator shares by the amount of voter baeth for the validator
			valAddrStr := valAddr.String()
			if val, ok := currValidators[valAddrStr]; ok {
				val.DelegatorDeductions = val.DelegatorDeductions.Add(coin.Amount.ToDec())
				currValidators[valAddrStr] = val
			}

			// votingPower = amount of uaeth coin
			stakedCoins, err := th.lk.GetStakedTokensForDerivatives(ctx, sdk.NewCoins(coin))
			if err != nil {
				// error is returned only if the baeth denom is incorrect, which should never happen here.
				panic(err)
			}
			votingPower := stakedCoins.Amount.ToDec()

			for _, option := range vote.Options {
				subPower := votingPower.Mul(option.Weight)
				results[option.Option] = results[option.Option].Add(subPower)
			}
			totalVotingPower = totalVotingPower.Add(votingPower)
		}

		th.gk.DeleteVote(ctx, vote.ProposalId, voter)
		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if len(val.Vote) == 0 {
			continue
		}

		sharesAfterDeductions := val.DelegatorShares.Sub(val.DelegatorDeductions)
		votingPower := sharesAfterDeductions.MulInt(val.BondedTokens).Quo(val.DelegatorShares)

		for _, option := range val.Vote {
			subPower := votingPower.Mul(option.Weight)
			results[option.Option] = results[option.Option].Add(subPower)
		}
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyParams := th.gk.GetTallyParams(ctx)
	tallyResults = types.NewTallyResultFromMap(results)

	// TODO: Upgrade the spec to cover all of these cases & remove pseudocode.
	// If there is no staked coins, the proposal fails
	if th.stk.TotalBondedTokens(ctx).IsZero() {
		return false, false, tallyResults
	}

	// If there is not enough quorum of votes, the proposal fails
	percentVoting := totalVotingPower.Quo(th.stk.TotalBondedTokens(ctx).ToDec())
	if percentVoting.LT(tallyParams.Quorum) {
		return false, true, tallyResults
	}

	// If no one votes (everyone abstains), proposal fails
	if totalVotingPower.Sub(results[types.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false, tallyResults
	}

	// If more than 1/3 of voters veto, proposal fails
	if results[types.OptionNoWithVeto].Quo(totalVotingPower).GT(tallyParams.VetoThreshold) {
		return false, true, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[types.OptionYes].Quo(totalVotingPower.Sub(results[types.OptionAbstain])).GT(tallyParams.Threshold) {
		return true, false, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails
	return false, false, tallyResults
}

// baethByDenom a map of the baeth denom and the amount of baeth for that denom.
type baethByDenom map[string]sdk.Int

func (baethMap baethByDenom) add(coin sdk.Coin) {
	_, found := baethMap[coin.Denom]
	if !found {
		baethMap[coin.Denom] = sdk.ZeroInt()
	}
	baethMap[coin.Denom] = baethMap[coin.Denom].Add(coin.Amount)
}

func (baethMap baethByDenom) toCoins() sdk.Coins {
	coins := sdk.Coins{}
	for denom, amt := range baethMap {
		coins = coins.Add(sdk.NewCoin(denom, amt))
	}
	return coins.Sort()
}

// getAddrBaeth returns a map of validator address & the amount of baeth
// of the addr for each validator.
func (th TallyHandler) getAddrBaeth(ctx sdk.Context, addr sdk.AccAddress) baethByDenom {
	results := make(baethByDenom)
	th.addBaethFromWallet(ctx, addr, results)
	th.addBaethFromSavings(ctx, addr, results)
	th.addBaethFromEarn(ctx, addr, results)
	return results
}

// addBaethFromWallet adds all addr balances of baeth in x/bank.
func (th TallyHandler) addBaethFromWallet(ctx sdk.Context, addr sdk.AccAddress, baeth baethByDenom) {
	coins := th.bk.GetAllBalances(ctx, addr)
	for _, coin := range coins {
		if th.lk.IsDerivativeDenom(ctx, coin.Denom) {
			baeth.add(coin)
		}
	}
}

// addBaethFromSavings adds all addr deposits of baeth in x/savings.
func (th TallyHandler) addBaethFromSavings(ctx sdk.Context, addr sdk.AccAddress, baeth baethByDenom) {
	deposit, found := th.svk.GetDeposit(ctx, addr)
	if !found {
		return
	}
	for _, coin := range deposit.Amount {
		if th.lk.IsDerivativeDenom(ctx, coin.Denom) {
			baeth.add(coin)
		}
	}
}

// addBaethFromEarn adds all addr deposits of baeth in x/earn.
func (th TallyHandler) addBaethFromEarn(ctx sdk.Context, addr sdk.AccAddress, baeth baethByDenom) {
	shares, found := th.ek.GetVaultAccountShares(ctx, addr)
	if !found {
		return
	}
	for _, share := range shares {
		if th.lk.IsDerivativeDenom(ctx, share.Denom) {
			if coin, err := th.ek.ConvertToAssets(ctx, share); err == nil {
				baeth.add(coin)
			}
		}
	}
}
