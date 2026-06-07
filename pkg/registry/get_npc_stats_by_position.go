package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetNPCStatsByPosition(positionZ consts.PlanetPositionZ) (NPCStats, error) {
	npcStats, ok := r.npcStats[positionZ]
	if !ok {
		return NPCStats{}, fmt.Errorf("%w: Position Z %d", ErrNotFound, positionZ)
	}

	return npcStats, nil
}
