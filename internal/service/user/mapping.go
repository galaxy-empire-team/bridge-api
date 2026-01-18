package user

import (
	"initialservice/internal/models"
	userstorage "initialservice/internal/storage/user"
)

func toStorageUser(u models.User) userstorage.User {
	return userstorage.User{
		ID:    u.ID,
		Login: u.Login,
	}
}

func toModelUser(u userstorage.User) models.User {
	return models.User{
		ID:    u.ID,
		Login: u.Login,
	}
}
