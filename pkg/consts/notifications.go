package consts

type NotificationID uint8

func (b NotificationID) ToUint8() uint8 {
	return uint8(b)
}

type NotificationType string

const (
	NotificationTypeColonize NotificationType = "colonize"
	NotificationTypeAttack   NotificationType = "attack"
	NotificationTypeReturn   NotificationType = "return"
)
