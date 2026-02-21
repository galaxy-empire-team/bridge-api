package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetMissionIDByType(missionType consts.MissionType) (consts.MissionID, error) {
	for id, mType := range r.missions {
		if mType == missionType {
			return id, nil
		}
	}

	return 0, fmt.Errorf("%w: type %s", ErrNotFound, missionType)
}
