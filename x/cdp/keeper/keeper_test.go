package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/mokitanetwork/aether/app"
	"github.com/mokitanetwork/aether/x/cdp/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	keeper keeper.Keeper
	app    app.TestApp
	ctx    sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.ResetChain()
}

func (suite *KeeperTestSuite) ResetChain() {
	tApp := app.NewTestApp()
	ctx := tApp.NewContext(true, tmproto.Header{Height: 1, Time: tmtime.Now()})
	keeper := tApp.GetCDPKeeper()

	suite.app = tApp
	suite.ctx = ctx
	suite.keeper = keeper
}
