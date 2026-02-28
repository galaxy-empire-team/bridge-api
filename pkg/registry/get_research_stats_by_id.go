package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetResearchStatsByID(researchID consts.ResearchID) (ResearchStats, error) {
	researchStats, exists := r.researches[researchID]
	if !exists {
		return ResearchStats{}, fmt.Errorf("%w: id %d", ErrNotFound, researchID)
	}

	return researchStats, nil
}
