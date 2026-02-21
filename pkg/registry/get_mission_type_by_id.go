package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetMissionTypeByID(missionID consts.MissionID) (consts.MissionType, error) {
	mType, exists := r.missions[missionID]
	if !exists {
		return "", fmt.Errorf("%w: id %d", ErrNotFound, missionID)
	}

	return mType, nil
}
