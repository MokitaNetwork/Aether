package types_test

import (
	"testing"

	"github.com/mokitanetwork/aether/x/swap/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	key := types.PoolKey(types.PoolID("uaeth", "usdx"))
	assert.Equal(t, types.PoolID("uaeth", "usdx"), string(key))

	key = types.DepositorPoolSharesKey(sdk.AccAddress("testaddress1"), types.PoolID("uaeth", "usdx"))
	assert.Equal(t, string(sdk.AccAddress("testaddress1"))+"|"+types.PoolID("uaeth", "usdx"), string(key))
}
