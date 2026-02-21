package registry

import (
	"fmt"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *Registry) GetNotificationTypeByID(notificationID consts.NotificationID) (consts.NotificationType, error) {
	nType, exists := r.notifications[notificationID]
	if !exists {
		return "", fmt.Errorf("%w: id %d", ErrNotFound, notificationID)
	}

	return nType, nil
}
