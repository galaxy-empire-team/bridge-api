package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetNotificationIDByType(notificationType consts.NotificationType) (consts.NotificationID, error) {
	for id, nType := range r.notifications {
		if nType == notificationType {
			return id, nil
		}
	}

	return 0, fmt.Errorf("%w: type %s", ErrNotFound, notificationType)
}
