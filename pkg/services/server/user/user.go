package user

import "mephiMainProject/pkg/services/server/config"

type UserRepo interface {
	Authorize(login, plainPassword string) (*config.UserAuthData, error)
	Register(login, plainPassword string) (int, error)
}
