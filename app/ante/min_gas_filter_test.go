package ante_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/mokitanetwork/aether/app"
	"github.com/mokitanetwork/aether/app/ante"
)

func mustParseDecCoins(value string) sdk.DecCoins {
	coins, err := sdk.ParseDecCoins(strings.ReplaceAll(value, ";", ","))
	if err != nil {
		panic(err)
	}

	return coins
}

func TestEvmMinGasFilter(t *testing.T) {
	tApp := app.NewTestApp()
	handler := ante.NewEvmMinGasFilter(tApp.GetEvmKeeper())

	ctx := tApp.NewContext(true, tmproto.Header{Height: 1, Time: tmtime.Now()})
	tApp.GetEvmKeeper().SetParams(ctx, evmtypes.Params{
		EvmDenom: "aaeth",
	})

	testCases := []struct {
		name                 string
		minGasPrices         sdk.DecCoins
		expectedMinGasPrices sdk.DecCoins
	}{
		{
			"no min gas prices",
			mustParseDecCoins(""),
			mustParseDecCoins(""),
		},
		{
			"zero uaeth gas price",
			mustParseDecCoins("0uaeth"),
			mustParseDecCoins("0uaeth"),
		},
		{
			"non-zero uaeth gas price",
			mustParseDecCoins("0.001uaeth"),
			mustParseDecCoins("0.001uaeth"),
		},
		{
			"zero uaeth gas price, min aaeth price",
			mustParseDecCoins("0uaeth;100000aaeth"),
			mustParseDecCoins("0uaeth"), // aaeth is removed
		},
		{
			"zero uaeth gas price, min aaeth price, other token",
			mustParseDecCoins("0uaeth;100000aaeth;0.001other"),
			mustParseDecCoins("0uaeth;0.001other"), // aaeth is removed
		},
		{
			"non-zero uaeth gas price, min aaeth price",
			mustParseDecCoins("0.25uaeth;100000aaeth;0.001other"),
			mustParseDecCoins("0.25uaeth;0.001other"), // aaeth is removed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tApp.NewContext(true, tmproto.Header{Height: 1, Time: tmtime.Now()})

			ctx = ctx.WithMinGasPrices(tc.minGasPrices)
			mmd := MockAnteHandler{}

			_, err := handler.AnteHandle(ctx, nil, false, mmd.AnteHandle)
			require.NoError(t, err)
			require.True(t, mmd.WasCalled)

			assert.NoError(t, mmd.CalledCtx.MinGasPrices().Validate())
			assert.Equal(t, tc.expectedMinGasPrices, mmd.CalledCtx.MinGasPrices())
		})
	}
}
