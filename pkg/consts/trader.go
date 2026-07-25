package consts

type TraderItemID uint8

type TraderItemType string

const (
	TraderItemTypeBoost          TraderItemType = "boost"
	TraderItemTypeSpaceship      TraderItemType = "spaceship"
	TraderItemTypeMoon           TraderItemType = "moon"
	TraderItemTypeAutoMistLaunch TraderItemType = "auto_mist_launch"
)
