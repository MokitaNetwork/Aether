package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/mokitanetwork/aether/x/evmutil/types"
)

const (
	// EvmDenom is the gas denom used by the evm
	EvmDenom = "aaeth"

	// CosmosDenom is the gas denom used by the aeth app
	CosmosDenom = "uaeth"
)

// ConversionMultiplier is the conversion multiplier between aaeth and uaeth
var ConversionMultiplier = sdk.NewInt(1_000_000_000_000)

var _ evmtypes.BankKeeper = EvmBankKeeper{}

// EvmBankKeeper is a BankKeeper wrapper for the x/evm module to allow the use
// of the 18 decimal aaeth coin on the evm.
// x/evm consumes gas and send coins by minting and burning aaeth coins in its module
// account and then sending the funds to the target account.
// This keeper uses both the uaeth coin and a separate aaeth balance to manage the
// extra percision needed by the evm.
type EvmBankKeeper struct {
	aaethKeeper Keeper
	bk          types.BankKeeper
	ak          types.AccountKeeper
}

func NewEvmBankKeeper(aaethKeeper Keeper, bk types.BankKeeper, ak types.AccountKeeper) EvmBankKeeper {
	return EvmBankKeeper{
		aaethKeeper: aaethKeeper,
		bk:          bk,
		ak:          ak,
	}
}

// GetBalance returns the total **spendable** balance of aaeth for a given account by address.
func (k EvmBankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	if denom != EvmDenom {
		panic(fmt.Errorf("only evm denom %s is supported by EvmBankKeeper", EvmDenom))
	}

	spendableCoins := k.bk.SpendableCoins(ctx, addr)
	uaeth := spendableCoins.AmountOf(CosmosDenom)
	aaeth := k.aaethKeeper.GetBalance(ctx, addr)
	total := uaeth.Mul(ConversionMultiplier).Add(aaeth)
	return sdk.NewCoin(EvmDenom, total)
}

// SendCoinsFromModuleToAccount transfers aaeth coins from a ModuleAccount to an AccAddress.
// It will panic if the module account does not exist. An error is returned if the recipient
// address is black-listed or if sending the tokens fails.
func (k EvmBankKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	uaeth, aaeth, err := SplitAaethCoins(amt)
	if err != nil {
		return err
	}

	if uaeth.Amount.IsPositive() {
		if err := k.bk.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, sdk.NewCoins(uaeth)); err != nil {
			return err
		}
	}

	senderAddr := k.GetModuleAddress(senderModule)
	if err := k.ConvertOneUaethToAaethIfNeeded(ctx, senderAddr, aaeth); err != nil {
		return err
	}

	if err := k.aaethKeeper.SendBalance(ctx, senderAddr, recipientAddr, aaeth); err != nil {
		return err
	}

	return k.ConvertAaethToUaeth(ctx, recipientAddr)
}

// SendCoinsFromAccountToModule transfers aaeth coins from an AccAddress to a ModuleAccount.
// It will panic if the module account does not exist.
func (k EvmBankKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	uaeth, aaethNeeded, err := SplitAaethCoins(amt)
	if err != nil {
		return err
	}

	if uaeth.IsPositive() {
		if err := k.bk.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, sdk.NewCoins(uaeth)); err != nil {
			return err
		}
	}

	if err := k.ConvertOneUaethToAaethIfNeeded(ctx, senderAddr, aaethNeeded); err != nil {
		return err
	}

	recipientAddr := k.GetModuleAddress(recipientModule)
	if err := k.aaethKeeper.SendBalance(ctx, senderAddr, recipientAddr, aaethNeeded); err != nil {
		return err
	}

	return k.ConvertAaethToUaeth(ctx, recipientAddr)
}

// MintCoins mints aaeth coins by minting the equivalent uaeth coins and any remaining aaeth coins.
// It will panic if the module account does not exist or is unauthorized.
func (k EvmBankKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	uaeth, aaeth, err := SplitAaethCoins(amt)
	if err != nil {
		return err
	}

	if uaeth.IsPositive() {
		if err := k.bk.MintCoins(ctx, moduleName, sdk.NewCoins(uaeth)); err != nil {
			return err
		}
	}

	recipientAddr := k.GetModuleAddress(moduleName)
	if err := k.aaethKeeper.AddBalance(ctx, recipientAddr, aaeth); err != nil {
		return err
	}

	return k.ConvertAaethToUaeth(ctx, recipientAddr)
}

// BurnCoins burns aaeth coins by burning the equivalent uaeth coins and any remaining aaeth coins.
// It will panic if the module account does not exist or is unauthorized.
func (k EvmBankKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	uaeth, aaeth, err := SplitAaethCoins(amt)
	if err != nil {
		return err
	}

	if uaeth.IsPositive() {
		if err := k.bk.BurnCoins(ctx, moduleName, sdk.NewCoins(uaeth)); err != nil {
			return err
		}
	}

	moduleAddr := k.GetModuleAddress(moduleName)
	if err := k.ConvertOneUaethToAaethIfNeeded(ctx, moduleAddr, aaeth); err != nil {
		return err
	}

	return k.aaethKeeper.RemoveBalance(ctx, moduleAddr, aaeth)
}

// ConvertOneUaethToAaethIfNeeded converts 1 uaeth to aaeth for an address if
// its aaeth balance is smaller than the aaethNeeded amount.
func (k EvmBankKeeper) ConvertOneUaethToAaethIfNeeded(ctx sdk.Context, addr sdk.AccAddress, aaethNeeded sdk.Int) error {
	aaethBal := k.aaethKeeper.GetBalance(ctx, addr)
	if aaethBal.GTE(aaethNeeded) {
		return nil
	}

	uaethToStore := sdk.NewCoins(sdk.NewCoin(CosmosDenom, sdk.OneInt()))
	if err := k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, uaethToStore); err != nil {
		return err
	}

	// add 1uaeth equivalent of aaeth to addr
	aaethToReceive := ConversionMultiplier
	if err := k.aaethKeeper.AddBalance(ctx, addr, aaethToReceive); err != nil {
		return err
	}

	return nil
}

// ConvertAaethToUaeth converts all available aaeth to uaeth for a given AccAddress.
func (k EvmBankKeeper) ConvertAaethToUaeth(ctx sdk.Context, addr sdk.AccAddress) error {
	totalAaeth := k.aaethKeeper.GetBalance(ctx, addr)
	uaeth, _, err := SplitAaethCoins(sdk.NewCoins(sdk.NewCoin(EvmDenom, totalAaeth)))
	if err != nil {
		return err
	}

	// do nothing if account does not have enough aaeth for a single uaeth
	uaethToReceive := uaeth.Amount
	if !uaethToReceive.IsPositive() {
		return nil
	}

	// remove aaeth used for converting to uaeth
	aaethToBurn := uaethToReceive.Mul(ConversionMultiplier)
	finalBal := totalAaeth.Sub(aaethToBurn)
	if err := k.aaethKeeper.SetBalance(ctx, addr, finalBal); err != nil {
		return err
	}

	fromAddr := k.GetModuleAddress(types.ModuleName)
	if err := k.bk.SendCoins(ctx, fromAddr, addr, sdk.NewCoins(uaeth)); err != nil {
		return err
	}

	return nil
}

func (k EvmBankKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	addr := k.ak.GetModuleAddress(moduleName)
	if addr == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
	}
	return addr
}

// SplitAaethCoins splits aaeth coins to the equivalent uaeth coins and any remaining aaeth balance.
// An error will be returned if the coins are not valid or if the coins are not the aaeth denom.
func SplitAaethCoins(coins sdk.Coins) (sdk.Coin, sdk.Int, error) {
	aaeth := sdk.ZeroInt()
	uaeth := sdk.NewCoin(CosmosDenom, sdk.ZeroInt())

	if len(coins) == 0 {
		return uaeth, aaeth, nil
	}

	if err := ValidateEvmCoins(coins); err != nil {
		return uaeth, aaeth, err
	}

	// note: we should always have len(coins) == 1 here since coins cannot have dup denoms after we validate.
	coin := coins[0]
	remainingBalance := coin.Amount.Mod(ConversionMultiplier)
	if remainingBalance.IsPositive() {
		aaeth = remainingBalance
	}
	uaethAmount := coin.Amount.Quo(ConversionMultiplier)
	if uaethAmount.IsPositive() {
		uaeth = sdk.NewCoin(CosmosDenom, uaethAmount)
	}

	return uaeth, aaeth, nil
}

// ValidateEvmCoins validates the coins from evm is valid and is the EvmDenom (aaeth).
func ValidateEvmCoins(coins sdk.Coins) error {
	if len(coins) == 0 {
		return nil
	}

	// validate that coins are non-negative, sorted, and no dup denoms
	if err := coins.Validate(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, coins.String())
	}

	// validate that coin denom is aaeth
	if len(coins) != 1 || coins[0].Denom != EvmDenom {
		errMsg := fmt.Sprintf("invalid evm coin denom, only %s is supported", EvmDenom)
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, errMsg)
	}

	return nil
}
