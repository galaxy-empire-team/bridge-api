package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetBoostStatsByTier(tier consts.BoostTier) (BoostStats, error) {
	for _, boostStats := range r.boosts {
		if boostStats.Tier == tier {
			return boostStats, nil
		}
	}

	return BoostStats{}, fmt.Errorf("%w: tier %d", ErrNotFound, tier)
}
