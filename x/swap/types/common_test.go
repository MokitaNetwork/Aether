package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mokitanetwork/aether/app"
)

func init() {
	aethConfig := sdk.GetConfig()
	app.SetBech32AddressPrefixes(aethConfig)
	app.SetBip44CoinType(aethConfig)
	aethConfig.Seal()
}
