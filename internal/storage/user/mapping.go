package user

import (
	"initialservice/internal/models"
)

func toStorageUser(u models.User) User {
	return User{
		ID:    u.ID,
		Login: u.Login,
	}
}

func toModelUser(u User) models.User {
	return models.User{
		ID:    u.ID,
		Login: u.Login,
	}
}
