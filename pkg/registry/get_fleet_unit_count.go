package registry

func (r *Registry) GetFleetUnitTypeCount() int {
	return len(r.fleet)
}
