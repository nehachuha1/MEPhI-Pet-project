package profile

import (
	"log"
	"mephiMainProject/pkg/services/server/config"
	"mephiMainProject/pkg/services/server/database"
)

type ProfileRepository struct {
	db *database.DatabaseORM
}

func NewProfileRepository(cfg *config.Config) *ProfileRepository {
	return &ProfileRepository{
		db: database.NewDBUsage(cfg),
	}
}

func (repo *ProfileRepository) CreateProfile(data *config.User) error {
	err := repo.db.CreateProfile(data, data.Login)
	if err != nil {
		log.Printf("ProfileRepo.CreateProfile err - %v", err)
		return err
	}
	return nil
}

func (repo *ProfileRepository) GetProfile(username string) (config.User, error) {
	currentProfile, err := repo.db.GetProfile(username)
	if err != nil {
		log.Printf("ProfileRepo.GetProfile err - %v", err)
		return config.User{}, err
	}
	return currentProfile, nil
}

func (repo *ProfileRepository) EditProfile(username string, newData *config.User) error {
	err := repo.db.EditProfile(username, newData)
	if err != nil {
		log.Printf("ProfileRepo.EditProfile err - %v", err)
		return err
	}
	return nil
}

func (repo *ProfileRepository) DeleteProfile(username string) error {
	err := repo.db.DeleteProfile(username)
	if err != nil {
		log.Printf("ProfileRepo.DeleteProfile err - %v", err)
		return err
	}
	return nil
}
