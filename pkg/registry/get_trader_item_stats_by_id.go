package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetTraderItemStatsByID(id consts.TraderItemID) (TraderItemStats, error) {
	traderStats, exists := r.traderStats[id]
	if !exists {
		return TraderItemStats{}, fmt.Errorf("%w: id %d", ErrNotFound, id)
	}

	return traderStats, nil
}
