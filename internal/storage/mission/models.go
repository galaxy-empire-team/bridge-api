package mission

type Fleet struct {
	Fleet []FleetUnit `json:"fleet"`
}

type FleetUnit struct {
	ID    uint8  `json:"id"`
	Count uint64 `json:"count"`
}
