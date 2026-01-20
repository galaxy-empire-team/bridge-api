package txmanager

import (
	"github.com/jackc/pgx/v5"

	planetstorage "initialservice/internal/storage/planet"
	userstorage "initialservice/internal/storage/user"
)

// I don't want to write boilerplate stuff, embed all storages ^_^
type StorageSet struct {
	*planetstorage.PlanetStorage
	*userstorage.UserStorage
}

func newStorageSet(tx pgx.Tx) StorageSet {
	return StorageSet{
		PlanetStorage: planetstorage.New(tx),
		UserStorage:   userstorage.New(tx),
	}
}
