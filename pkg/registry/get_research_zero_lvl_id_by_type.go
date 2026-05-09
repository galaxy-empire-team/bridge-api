package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetResearchZeroLvlIDByType(researchType consts.ResearchType) (consts.ResearchID, error) {
	researchID, exists := r.zeroLvlResearches[researchType]
	if !exists {
		return 0, fmt.Errorf("%w: type %s", ErrNotFound, researchType)
	}

	return researchID, nil
}
