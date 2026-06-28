package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetBoostStatsByID(id consts.BoostID) (BoostStats, error) {
	boostStats, exists := r.boosts[id]
	if !exists {
		return BoostStats{}, fmt.Errorf("%w: id %d", ErrNotFound, id)
	}

	return boostStats, nil
}
