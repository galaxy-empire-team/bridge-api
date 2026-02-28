package mission

type Fleet struct {
	Fleet []FleetUnit `json:"fleet"`
}

type FleetUnit struct {
	ID    uint8  `json:"id"`
	Count uint64 `json:"count"`
}

type Resources struct {
	Metal   uint64 `json:"metal"`
	Crystal uint64 `json:"crystal"`
	Gas     uint64 `json:"gas"`
}
