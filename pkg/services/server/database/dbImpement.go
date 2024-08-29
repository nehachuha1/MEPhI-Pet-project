package database

import "mephiMainProject/pkg/services/server/config"

type DatabaseControl interface {
	GetUserByID(id int) (config.User, error)
	GetUserByLogin(login string) (config.User, error)
	GetAuthUserData(login string) (*config.UserAuthData, error)
	RegisterUser(login, plainPassword string) (int, error)

	CheckSession(in *config.SessionID) (*config.Session, error)
	CreateSession(in *config.Session) (*config.Session, error) // в in не содержится айди сессии, в функции его присваиваем
	DeleteSession(in *config.SessionID) error
}
