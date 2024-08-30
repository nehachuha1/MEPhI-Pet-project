package user

import (
	"errors"
	argonpass "github.com/dwin/goArgonPass"
	"log"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/database"
)

var (
	ErrExistUser = errors.New("This user already exists")
	ErrNoUser    = errors.New("There's no user")
	ErrBadPass   = errors.New("Wrong password")
)

type UserRepository struct {
	dbORM *database.DatabaseORM
}

func NewUserRepository(cfg *config.Config) *UserRepository {
	return &UserRepository{
		dbORM: database.NewDBUsage(cfg),
	}
}

func (repo *UserRepository) Authorize(login, plainPassword string) (*config.UserAuthData, error) {
	currentUserAuthData, err := repo.dbORM.GetAuthUserData(login)
	if err != nil {
		log.Printf("From UserRepo to Auth table err - %v", err)
		return &config.UserAuthData{}, err
	}
	isPasswordMatch := argonpass.Verify(plainPassword, currentUserAuthData.Password)
	if isPasswordMatch != nil {
		log.Printf("User passwords dont match")
		return &config.UserAuthData{}, isPasswordMatch
	}
	return currentUserAuthData, nil
}

func (repo *UserRepository) Register(login, plainPassword string) (int, error) {
	lastID, err := repo.dbORM.RegisterUser(login, plainPassword)
	if err != nil {
		log.Printf("Register user from repo err - %v", err)
		return 0, err
	}
	return lastID, nil
}
