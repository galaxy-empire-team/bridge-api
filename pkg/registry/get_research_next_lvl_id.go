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

	// For simplicity, I assume that the next level of a research has an ID that is exactly 1 greater than the current research's ID.
	// TODO rewrite using map and linked list
	nextLvlResearch, ok := r.researches[researchID+1]
	if !ok {
		return 0, ErrMaxLevelReached
	}

	if currentResearch.Type != nextLvlResearch.Type {
		return 0, ErrMaxLevelReached
	}

	return researchID + 1, nil
}
