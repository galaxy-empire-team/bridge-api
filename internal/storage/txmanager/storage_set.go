package txmanager

import (
	"github.com/jackc/pgx/v5"

	missionstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/mission"
	planetstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/planet"
	researchstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/research"
)

// I don't want to write boilerplate stuff, embed all storages ^_^
type planetMissionStorageSet struct {
	*planetstorage.PlanetStorage
	*missionstorage.MissionStorage
}

func newPlanetMissionStorageSet(tx pgx.Tx) planetMissionStorageSet {
	return planetMissionStorageSet{
		PlanetStorage:  planetstorage.New(tx),
		MissionStorage: missionstorage.New(tx),
	}
}

type planetResearchStorageSet struct {
	*planetstorage.PlanetStorage
	*researchstorage.ResearchStorage
}

func newPlanetResearchStorageSet(tx pgx.Tx) planetResearchStorageSet {
	return planetResearchStorageSet{
		PlanetStorage:   planetstorage.New(tx),
		ResearchStorage: researchstorage.New(tx),
	}
}
