package txmanager

import (
	"github.com/jackc/pgx/v5"

	missionstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/mission"
	planetstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/planet"
	researchstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/research"
	userstorage "github.com/galaxy-empire-team/bridge-api/internal/storage/user"
)

// I don't want to write boilerplate stuff, embed all storages ^_^
type StorageSet struct {
	*planetstorage.PlanetStorage
	*userstorage.UserStorage
	*missionstorage.MissionStorage
	*researchstorage.ResearchStorage
}

func newStorageSet(tx pgx.Tx) StorageSet {
	return StorageSet{
		PlanetStorage:   planetstorage.New(tx),
		UserStorage:     userstorage.New(tx),
		MissionStorage:  missionstorage.New(tx),
		ResearchStorage: researchstorage.New(tx),
	}
}
