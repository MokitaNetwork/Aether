package v0_16

import (
	v015aethdist "github.com/mokitanetwork/aether/x/aethdist/legacy/v0_15"
	v016aethdist "github.com/mokitanetwork/aether/x/aethdist/types"
)

func migrateParams(oldParams v015aethdist.Params) v016aethdist.Params {
	periods := make([]v016aethdist.Period, len(oldParams.Periods))
	for i, oldPeriod := range oldParams.Periods {
		periods[i] = v016aethdist.Period{
			Start:     oldPeriod.Start,
			End:       oldPeriod.End,
			Inflation: oldPeriod.Inflation,
		}
	}
	return v016aethdist.Params{
		Periods: periods,
		Active:  oldParams.Active,
	}
}

// Migrate converts v0.15 aethdist state and returns it in v0.16 format
func Migrate(oldState v015aethdist.GenesisState) *v016aethdist.GenesisState {
	return &v016aethdist.GenesisState{
		Params:            migrateParams(oldState.Params),
		PreviousBlockTime: oldState.PreviousBlockTime,
	}
}
