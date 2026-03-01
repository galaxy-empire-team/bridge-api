package registry

func (r *Registry) GetAllBuildingStats() ([]BuildingStats, error) {
	stats := make([]BuildingStats, 0, len(r.buildings))
	for _, stat := range r.buildings {
		stats = append(stats, stat)
	}

	return stats, nil
}
