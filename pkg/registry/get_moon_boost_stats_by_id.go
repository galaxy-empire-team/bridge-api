package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetMoonBoostStatsByID(id consts.MoonBoostID) (MoonBoostStats, error) {
	moonBoostStats, exists := r.moonBoosts[id]
	if !exists {
		return MoonBoostStats{}, fmt.Errorf("%w: id %d", ErrNotFound, id)
	}

	return moonBoostStats, nil
}
