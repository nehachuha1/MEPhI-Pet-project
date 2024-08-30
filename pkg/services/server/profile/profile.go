package profile

import "mephiMainProject/pkg/services/server/config"

type ProfileRepo interface {
	CreateProfile(data *config.User, username string) error
	GetProfile(username string) (config.User, error)
	EditProfile(username string, newData *config.User) error
	DeleteProfile(username string) error
}
