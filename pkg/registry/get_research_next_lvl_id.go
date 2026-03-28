package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetResearchNextLvlID(researchID consts.ResearchID) (consts.ResearchID, error) {
	currentResearch, ok := r.researches[researchID]
	if !ok {
		return 0, fmt.Errorf("%w: ID %d", ErrNotFound, researchID)
	}

	if currentResearch.Type == consts.ResearchTypeColonizeTechnology {
		if currentResearch.Level >= consts.MaxColonizationResearchLvl {
			return 0, ErrMaxLevelReached
		}
	}

	if currentResearch.Level >= consts.MaxResearchLvl {
		return 0, ErrMaxLevelReached
	}

	// For simplicity, I assume that the next level of a research has an ID that is exactly 1 greater than the current research's ID.
	// TODO rewrite using map and linked list
	_, ok = r.researches[researchID+1]
	if !ok {
		return 0, fmt.Errorf("%w: ID %d", ErrNotFound, researchID+1)
	}

	return researchID + 1, nil
}
