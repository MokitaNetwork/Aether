package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmtime "github.com/tendermint/tendermint/types/time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/mokitanetwork/aether/x/evmutil/keeper"
	"github.com/mokitanetwork/aether/x/evmutil/testutil"
	"github.com/mokitanetwork/aether/x/evmutil/types"
)

type evmBankKeeperTestSuite struct {
	testutil.Suite
}

func (suite *evmBankKeeperTestSuite) SetupTest() {
	suite.Suite.SetupTest()
}

func (suite *evmBankKeeperTestSuite) TestGetBalance_ReturnsSpendable() {
	startingCoins := sdk.NewCoins(sdk.NewInt64Coin("uaeth", 10))
	startingAaeth := sdk.NewInt(100)

	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)
	bacc := authtypes.NewBaseAccountWithAddress(suite.Addrs[0])
	vacc := vesting.NewContinuousVestingAccount(bacc, startingCoins, now.Unix(), endTime.Unix())
	suite.AccountKeeper.SetAccount(suite.Ctx, vacc)

	err := suite.App.FundAccount(suite.Ctx, suite.Addrs[0], startingCoins)
	suite.Require().NoError(err)
	err = suite.Keeper.SetBalance(suite.Ctx, suite.Addrs[0], startingAaeth)
	suite.Require().NoError(err)

	coin := suite.EvmBankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "aaeth")
	suite.Require().Equal(startingAaeth, coin.Amount)

	ctx := suite.Ctx.WithBlockTime(now.Add(12 * time.Hour))
	coin = suite.EvmBankKeeper.GetBalance(ctx, suite.Addrs[0], "aaeth")
	suite.Require().Equal(sdk.NewIntFromUint64(5_000_000_000_100), coin.Amount)
}

func (suite *evmBankKeeperTestSuite) TestGetBalance_NotEvmDenom() {
	suite.Require().Panics(func() {
		suite.EvmBankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "uaeth")
	})
	suite.Require().Panics(func() {
		suite.EvmBankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "busd")
	})
}

func (suite *evmBankKeeperTestSuite) TestGetBalance() {
	tests := []struct {
		name           string
		startingAmount sdk.Coins
		expAmount      sdk.Int
	}{
		{
			"uaeth with aaeth",
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 100),
				sdk.NewInt64Coin("uaeth", 10),
			),
			sdk.NewInt(10_000_000_000_100),
		},
		{
			"just aaeth",
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 100),
				sdk.NewInt64Coin("busd", 100),
			),
			sdk.NewInt(100),
		},
		{
			"just uaeth",
			sdk.NewCoins(
				sdk.NewInt64Coin("uaeth", 10),
				sdk.NewInt64Coin("busd", 100),
			),
			sdk.NewInt(10_000_000_000_000),
		},
		{
			"no uaeth or aaeth",
			sdk.NewCoins(),
			sdk.ZeroInt(),
		},
		{
			"with avaka that is more than 1 uaeth",
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 20_000_000_000_220),
				sdk.NewInt64Coin("uaeth", 11),
			),
			sdk.NewInt(31_000_000_000_220),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.FundAccountWithAether(suite.Addrs[0], tt.startingAmount)
			coin := suite.EvmBankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "aaeth")
			suite.Require().Equal(tt.expAmount, coin.Amount)
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestSendCoinsFromModuleToAccount() {
	startingModuleCoins := sdk.NewCoins(
		sdk.NewInt64Coin("aaeth", 200),
		sdk.NewInt64Coin("uaeth", 100),
	)
	tests := []struct {
		name           string
		sendCoins      sdk.Coins
		startingAccBal sdk.Coins
		expAccBal      sdk.Coins
		hasErr         bool
	}{
		{
			"send more than 1 uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 12_000_000_000_010)),
			sdk.Coins{},
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 10),
				sdk.NewInt64Coin("uaeth", 12),
			),
			false,
		},
		{
			"send less than 1 uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 122)),
			sdk.Coins{},
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 122),
				sdk.NewInt64Coin("uaeth", 0),
			),
			false,
		},
		{
			"send an exact amount of uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 98_000_000_000_000)),
			sdk.Coins{},
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 0o0),
				sdk.NewInt64Coin("uaeth", 98),
			),
			false,
		},
		{
			"send no aaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 0)),
			sdk.Coins{},
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 0),
				sdk.NewInt64Coin("uaeth", 0),
			),
			false,
		},
		{
			"errors if sending other coins",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 500), sdk.NewInt64Coin("busd", 1000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"errors if not enough total aaeth to cover",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100_000_000_001_000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"errors if not enough uaeth to cover",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 200_000_000_000_000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"converts receiver's aaeth to uaeth if there's enough aaeth after the transfer",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 99_000_000_000_200)),
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 999_999_999_900),
				sdk.NewInt64Coin("uaeth", 1),
			),
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 100),
				sdk.NewInt64Coin("uaeth", 101),
			),
			false,
		},
		{
			"converts all of receiver's aaeth to uaeth even if somehow receiver has more than 1uaeth of aaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 12_000_000_000_100)),
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 5_999_999_999_990),
				sdk.NewInt64Coin("uaeth", 1),
			),
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 90),
				sdk.NewInt64Coin("uaeth", 19),
			),
			false,
		},
		{
			"swap 1 uaeth for aaeth if module account doesn't have enough aaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 99_000_000_001_000)),
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 200),
				sdk.NewInt64Coin("uaeth", 1),
			),
			sdk.NewCoins(
				sdk.NewInt64Coin("aaeth", 1200),
				sdk.NewInt64Coin("uaeth", 100),
			),
			false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.FundAccountWithAether(suite.Addrs[0], tt.startingAccBal)
			suite.FundModuleAccountWithAether(evmtypes.ModuleName, startingModuleCoins)

			// fund our module with some uaeth to account for converting extra aaeth back to uaeth
			suite.FundModuleAccountWithAether(types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("uaeth", 10)))

			err := suite.EvmBankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, evmtypes.ModuleName, suite.Addrs[0], tt.sendCoins)
			if tt.hasErr {
				suite.Require().Error(err)
				return
			} else {
				suite.Require().NoError(err)
			}

			// check uaeth
			uaethSender := suite.BankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "uaeth")
			suite.Require().Equal(tt.expAccBal.AmountOf("uaeth").Int64(), uaethSender.Amount.Int64())

			// check aaeth
			actualAaeth := suite.Keeper.GetBalance(suite.Ctx, suite.Addrs[0])
			suite.Require().Equal(tt.expAccBal.AmountOf("aaeth").Int64(), actualAaeth.Int64())
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestSendCoinsFromAccountToModule() {
	startingAccCoins := sdk.NewCoins(
		sdk.NewInt64Coin("aaeth", 200),
		sdk.NewInt64Coin("uaeth", 100),
	)
	startingModuleCoins := sdk.NewCoins(
		sdk.NewInt64Coin("aaeth", 100_000_000_000),
	)
	tests := []struct {
		name           string
		sendCoins      sdk.Coins
		expSenderCoins sdk.Coins
		expModuleCoins sdk.Coins
		hasErr         bool
	}{
		{
			"send more than 1 uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 12_000_000_000_010)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 190), sdk.NewInt64Coin("uaeth", 88)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100_000_000_010), sdk.NewInt64Coin("uaeth", 12)),
			false,
		},
		{
			"send less than 1 uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 122)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 78), sdk.NewInt64Coin("uaeth", 100)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100_000_000_122), sdk.NewInt64Coin("uaeth", 0)),
			false,
		},
		{
			"send an exact amount of uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 98_000_000_000_000)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 200), sdk.NewInt64Coin("uaeth", 2)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100_000_000_000), sdk.NewInt64Coin("uaeth", 98)),
			false,
		},
		{
			"send no aaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 0)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 200), sdk.NewInt64Coin("uaeth", 100)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100_000_000_000), sdk.NewInt64Coin("uaeth", 0)),
			false,
		},
		{
			"errors if sending other coins",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 500), sdk.NewInt64Coin("busd", 1000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"errors if have dup coins",
			sdk.Coins{
				sdk.NewInt64Coin("aaeth", 12_000_000_000_000),
				sdk.NewInt64Coin("aaeth", 2_000_000_000_000),
			},
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"errors if not enough total aaeth to cover",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100_000_000_001_000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"errors if not enough uaeth to cover",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 200_000_000_000_000)),
			sdk.Coins{},
			sdk.Coins{},
			true,
		},
		{
			"converts 1 uaeth to aaeth if not enough aaeth to cover",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 99_001_000_000_000)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 999_000_000_200), sdk.NewInt64Coin("uaeth", 0)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 101_000_000_000), sdk.NewInt64Coin("uaeth", 99)),
			false,
		},
		{
			"converts receiver's aaeth to uaeth if there's enough aaeth after the transfer",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 5_900_000_000_200)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100_000_000_000), sdk.NewInt64Coin("uaeth", 94)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 200), sdk.NewInt64Coin("uaeth", 6)),
			false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.FundAccountWithAether(suite.Addrs[0], startingAccCoins)
			suite.FundModuleAccountWithAether(evmtypes.ModuleName, startingModuleCoins)

			err := suite.EvmBankKeeper.SendCoinsFromAccountToModule(suite.Ctx, suite.Addrs[0], evmtypes.ModuleName, tt.sendCoins)
			if tt.hasErr {
				suite.Require().Error(err)
				return
			} else {
				suite.Require().NoError(err)
			}

			// check sender balance
			uaethSender := suite.BankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "uaeth")
			suite.Require().Equal(tt.expSenderCoins.AmountOf("uaeth").Int64(), uaethSender.Amount.Int64())
			actualAaeth := suite.Keeper.GetBalance(suite.Ctx, suite.Addrs[0])
			suite.Require().Equal(tt.expSenderCoins.AmountOf("aaeth").Int64(), actualAaeth.Int64())

			// check module balance
			moduleAddr := suite.AccountKeeper.GetModuleAddress(evmtypes.ModuleName)
			uaethSender = suite.BankKeeper.GetBalance(suite.Ctx, moduleAddr, "uaeth")
			suite.Require().Equal(tt.expModuleCoins.AmountOf("uaeth").Int64(), uaethSender.Amount.Int64())
			actualAaeth = suite.Keeper.GetBalance(suite.Ctx, moduleAddr)
			suite.Require().Equal(tt.expModuleCoins.AmountOf("aaeth").Int64(), actualAaeth.Int64())
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestBurnCoins() {
	startingUaeth := sdk.NewInt(100)
	tests := []struct {
		name       string
		burnCoins  sdk.Coins
		expUaeth   sdk.Int
		expAaeth   sdk.Int
		hasErr     bool
		aaethStart sdk.Int
	}{
		{
			"burn more than 1 uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 12_021_000_000_002)),
			sdk.NewInt(88),
			sdk.NewInt(100_000_000_000),
			false,
			sdk.NewInt(121_000_000_002),
		},
		{
			"burn less than 1 uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 122)),
			sdk.NewInt(100),
			sdk.NewInt(878),
			false,
			sdk.NewInt(1000),
		},
		{
			"burn an exact amount of uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 98_000_000_000_000)),
			sdk.NewInt(2),
			sdk.NewInt(10),
			false,
			sdk.NewInt(10),
		},
		{
			"burn no aaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 0)),
			startingUaeth,
			sdk.ZeroInt(),
			false,
			sdk.ZeroInt(),
		},
		{
			"errors if burning other coins",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 500), sdk.NewInt64Coin("busd", 1000)),
			startingUaeth,
			sdk.NewInt(100),
			true,
			sdk.NewInt(100),
		},
		{
			"errors if have dup coins",
			sdk.Coins{
				sdk.NewInt64Coin("aaeth", 12_000_000_000_000),
				sdk.NewInt64Coin("aaeth", 2_000_000_000_000),
			},
			startingUaeth,
			sdk.ZeroInt(),
			true,
			sdk.ZeroInt(),
		},
		{
			"errors if burn amount is negative",
			sdk.Coins{sdk.Coin{Denom: "aaeth", Amount: sdk.NewInt(-100)}},
			startingUaeth,
			sdk.NewInt(50),
			true,
			sdk.NewInt(50),
		},
		{
			"errors if not enough aaeth to cover burn",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100_999_000_000_000)),
			sdk.NewInt(0),
			sdk.NewInt(99_000_000_000),
			true,
			sdk.NewInt(99_000_000_000),
		},
		{
			"errors if not enough uaeth to cover burn",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 200_000_000_000_000)),
			sdk.NewInt(100),
			sdk.ZeroInt(),
			true,
			sdk.ZeroInt(),
		},
		{
			"converts 1 uaeth to aaeth if not enough aaeth to cover",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 12_021_000_000_002)),
			sdk.NewInt(87),
			sdk.NewInt(980_000_000_000),
			false,
			sdk.NewInt(1_000_000_002),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			startingCoins := sdk.NewCoins(
				sdk.NewCoin("uaeth", startingUaeth),
				sdk.NewCoin("aaeth", tt.aaethStart),
			)
			suite.FundModuleAccountWithAether(evmtypes.ModuleName, startingCoins)

			err := suite.EvmBankKeeper.BurnCoins(suite.Ctx, evmtypes.ModuleName, tt.burnCoins)
			if tt.hasErr {
				suite.Require().Error(err)
				return
			} else {
				suite.Require().NoError(err)
			}

			// check uaeth
			uaethActual := suite.BankKeeper.GetBalance(suite.Ctx, suite.EvmModuleAddr, "uaeth")
			suite.Require().Equal(tt.expUaeth, uaethActual.Amount)

			// check aaeth
			aaethActual := suite.Keeper.GetBalance(suite.Ctx, suite.EvmModuleAddr)
			suite.Require().Equal(tt.expAaeth, aaethActual)
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestMintCoins() {
	tests := []struct {
		name       string
		mintCoins  sdk.Coins
		uaeth      sdk.Int
		aaeth      sdk.Int
		hasErr     bool
		aaethStart sdk.Int
	}{
		{
			"mint more than 1 uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 12_021_000_000_002)),
			sdk.NewInt(12),
			sdk.NewInt(21_000_000_002),
			false,
			sdk.ZeroInt(),
		},
		{
			"mint less than 1 uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 901_000_000_001)),
			sdk.ZeroInt(),
			sdk.NewInt(901_000_000_001),
			false,
			sdk.ZeroInt(),
		},
		{
			"mint an exact amount of uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 123_000_000_000_000_000)),
			sdk.NewInt(123_000),
			sdk.ZeroInt(),
			false,
			sdk.ZeroInt(),
		},
		{
			"mint no aaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 0)),
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			false,
			sdk.ZeroInt(),
		},
		{
			"errors if minting other coins",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 500), sdk.NewInt64Coin("busd", 1000)),
			sdk.ZeroInt(),
			sdk.NewInt(100),
			true,
			sdk.NewInt(100),
		},
		{
			"errors if have dup coins",
			sdk.Coins{
				sdk.NewInt64Coin("aaeth", 12_000_000_000_000),
				sdk.NewInt64Coin("aaeth", 2_000_000_000_000),
			},
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			true,
			sdk.ZeroInt(),
		},
		{
			"errors if mint amount is negative",
			sdk.Coins{sdk.Coin{Denom: "aaeth", Amount: sdk.NewInt(-100)}},
			sdk.ZeroInt(),
			sdk.NewInt(50),
			true,
			sdk.NewInt(50),
		},
		{
			"adds to existing aaeth balance",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 12_021_000_000_002)),
			sdk.NewInt(12),
			sdk.NewInt(21_000_000_102),
			false,
			sdk.NewInt(100),
		},
		{
			"convert aaeth balance to uaeth if it exceeds 1 uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 10_999_000_000_000)),
			sdk.NewInt(12),
			sdk.NewInt(1_200_000_001),
			false,
			sdk.NewInt(1_002_200_000_001),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.FundModuleAccountWithAether(types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("uaeth", 10)))
			suite.FundModuleAccountWithAether(evmtypes.ModuleName, sdk.NewCoins(sdk.NewCoin("aaeth", tt.aaethStart)))

			err := suite.EvmBankKeeper.MintCoins(suite.Ctx, evmtypes.ModuleName, tt.mintCoins)
			if tt.hasErr {
				suite.Require().Error(err)
				return
			} else {
				suite.Require().NoError(err)
			}

			// check uaeth
			uaethActual := suite.BankKeeper.GetBalance(suite.Ctx, suite.EvmModuleAddr, "uaeth")
			suite.Require().Equal(tt.uaeth, uaethActual.Amount)

			// check aaeth
			aaethActual := suite.Keeper.GetBalance(suite.Ctx, suite.EvmModuleAddr)
			suite.Require().Equal(tt.aaeth, aaethActual)
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestValidateEvmCoins() {
	tests := []struct {
		name      string
		coins     sdk.Coins
		shouldErr bool
	}{
		{
			"valid coins",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 500)),
			false,
		},
		{
			"dup coins",
			sdk.Coins{sdk.NewInt64Coin("aaeth", 500), sdk.NewInt64Coin("aaeth", 500)},
			true,
		},
		{
			"not evm coins",
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 500)),
			true,
		},
		{
			"negative coins",
			sdk.Coins{sdk.Coin{Denom: "aaeth", Amount: sdk.NewInt(-500)}},
			true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := keeper.ValidateEvmCoins(tt.coins)
			if tt.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestConvertOneUaethToAaethIfNeeded() {
	aaethNeeded := sdk.NewInt(200)
	tests := []struct {
		name          string
		startingCoins sdk.Coins
		expectedCoins sdk.Coins
		success       bool
	}{
		{
			"not enough uaeth for conversion",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100)),
			false,
		},
		{
			"converts 1 uaeth to aaeth",
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 10), sdk.NewInt64Coin("aaeth", 100)),
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 9), sdk.NewInt64Coin("aaeth", 1_000_000_000_100)),
			true,
		},
		{
			"conversion not needed",
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 10), sdk.NewInt64Coin("aaeth", 200)),
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 10), sdk.NewInt64Coin("aaeth", 200)),
			true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			suite.FundAccountWithAether(suite.Addrs[0], tt.startingCoins)
			err := suite.EvmBankKeeper.ConvertOneUaethToAaethIfNeeded(suite.Ctx, suite.Addrs[0], aaethNeeded)
			moduleAether := suite.BankKeeper.GetBalance(suite.Ctx, suite.AccountKeeper.GetModuleAddress(types.ModuleName), "uaeth")
			if tt.success {
				suite.Require().NoError(err)
				if tt.startingCoins.AmountOf("aaeth").LT(aaethNeeded) {
					suite.Require().Equal(sdk.OneInt(), moduleAether.Amount)
				}
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(sdk.ZeroInt(), moduleAether.Amount)
			}

			aaeth := suite.Keeper.GetBalance(suite.Ctx, suite.Addrs[0])
			suite.Require().Equal(tt.expectedCoins.AmountOf("aaeth"), aaeth)
			uaeth := suite.BankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "uaeth")
			suite.Require().Equal(tt.expectedCoins.AmountOf("uaeth"), uaeth.Amount)
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestConvertAaethToUaeth() {
	tests := []struct {
		name          string
		startingCoins sdk.Coins
		expectedCoins sdk.Coins
	}{
		{
			"not enough uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 100), sdk.NewInt64Coin("uaeth", 0)),
		},
		{
			"converts aaeth for 1 uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 10), sdk.NewInt64Coin("aaeth", 1_000_000_000_003)),
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 11), sdk.NewInt64Coin("aaeth", 3)),
		},
		{
			"converts more than 1 uaeth of aaeth",
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 10), sdk.NewInt64Coin("aaeth", 8_000_000_000_123)),
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 18), sdk.NewInt64Coin("aaeth", 123)),
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.SetupTest()

			err := suite.App.FundModuleAccount(suite.Ctx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin("uaeth", 10)))
			suite.Require().NoError(err)
			suite.FundAccountWithAether(suite.Addrs[0], tt.startingCoins)
			err = suite.EvmBankKeeper.ConvertAaethToUaeth(suite.Ctx, suite.Addrs[0])
			suite.Require().NoError(err)
			aaeth := suite.Keeper.GetBalance(suite.Ctx, suite.Addrs[0])
			suite.Require().Equal(tt.expectedCoins.AmountOf("aaeth"), aaeth)
			uaeth := suite.BankKeeper.GetBalance(suite.Ctx, suite.Addrs[0], "uaeth")
			suite.Require().Equal(tt.expectedCoins.AmountOf("uaeth"), uaeth.Amount)
		})
	}
}

func (suite *evmBankKeeperTestSuite) TestSplitAaethCoins() {
	tests := []struct {
		name          string
		coins         sdk.Coins
		expectedCoins sdk.Coins
		shouldErr     bool
	}{
		{
			"invalid coins",
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 500)),
			nil,
			true,
		},
		{
			"empty coins",
			sdk.NewCoins(),
			sdk.NewCoins(),
			false,
		},
		{
			"uaeth & aaeth coins",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 8_000_000_000_123)),
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 8), sdk.NewInt64Coin("aaeth", 123)),
			false,
		},
		{
			"only aaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 10_123)),
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 10_123)),
			false,
		},
		{
			"only uaeth",
			sdk.NewCoins(sdk.NewInt64Coin("aaeth", 5_000_000_000_000)),
			sdk.NewCoins(sdk.NewInt64Coin("uaeth", 5)),
			false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			uaeth, aaeth, err := keeper.SplitAaethCoins(tt.coins)
			if tt.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tt.expectedCoins.AmountOf("uaeth"), uaeth.Amount)
				suite.Require().Equal(tt.expectedCoins.AmountOf("aaeth"), aaeth)
			}
		})
	}
}

func TestEvmBankKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(evmBankKeeperTestSuite))
}
