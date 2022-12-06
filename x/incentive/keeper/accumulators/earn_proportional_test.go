package accumulators_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mokitanetwork/aether/x/incentive/keeper/accumulators"
	"github.com/mokitanetwork/aether/x/incentive/types"
	"github.com/stretchr/testify/require"
)

func TestGetProportionalRewardPeriod(t *testing.T) {
	tests := []struct {
		name                  string
		giveRewardPeriod      types.MultiRewardPeriod
		giveTotalBaethSupply  sdk.Int
		giveSingleBaethSupply sdk.Int
		wantRewardsPerSecond  sdk.DecCoins
	}{
		{
			"full amount",
			types.NewMultiRewardPeriod(
				true,
				"",
				time.Time{},
				time.Time{},
				cs(c("uaeth", 100), c("hard", 200)),
			),
			i(100),
			i(100),
			toDcs(c("uaeth", 100), c("hard", 200)),
		},
		{
			"3/4 amount",
			types.NewMultiRewardPeriod(
				true,
				"",
				time.Time{},
				time.Time{},
				cs(c("uaeth", 100), c("hard", 200)),
			),
			i(10_000000),
			i(7_500000),
			toDcs(c("uaeth", 75), c("hard", 150)),
		},
		{
			"half amount",
			types.NewMultiRewardPeriod(
				true,
				"",
				time.Time{},
				time.Time{},
				cs(c("uaeth", 100), c("hard", 200)),
			),
			i(100),
			i(50),
			toDcs(c("uaeth", 50), c("hard", 100)),
		},
		{
			"under 1 unit",
			types.NewMultiRewardPeriod(
				true,
				"",
				time.Time{},
				time.Time{},
				cs(c("uaeth", 100), c("hard", 200)),
			),
			i(1000), // total baeth
			i(1),    // baeth supply of this specific vault
			dcs(dc("uaeth", "0.1"), dc("hard", "0.2")), // rewards per second rounded to 0 if under 1uaeth/1hard
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewardsPerSecond := accumulators.GetProportionalRewardsPerSecond(
				tt.giveRewardPeriod,
				tt.giveTotalBaethSupply,
				tt.giveSingleBaethSupply,
			)

			require.Equal(t, tt.wantRewardsPerSecond, rewardsPerSecond)
		})
	}
}
