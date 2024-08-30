package database

import "mephiMainProject/pkg/services/server/config"

type DatabaseControl interface {
	GetAuthUserData(login string) (*config.UserAuthData, error)
	RegisterUser(login, plainPassword string) (int, error)

	CheckSession(in *config.SessionID) (*config.Session, error)
	CreateSession(in *config.Session) (*config.Session, error) // в in не содержится айди сессии, в функции его присваиваем
	DeleteSession(in *config.SessionID) error

	CreateProfile(in *config.User, username string) error
	GetProfile(username string) (config.User, error)
	EditProfile(username string, newData *config.User) error
	DeleteProfile(username string) error
}
