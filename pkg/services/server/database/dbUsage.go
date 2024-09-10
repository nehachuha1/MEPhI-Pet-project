package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	argonpass "github.com/dwin/goArgonPass"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"math/rand"
	"mephiMainProject/pkg/services/server/config"
)

var (
	letterRunes            = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	sessionIDKeyLen        = 32
	ErrorNewSession        = errors.New("error by creating new session")
	ErrorResultNotOK       = errors.New("error: result in not OK")
	ErrorUserAlreadyExists = errors.New("user already exists")
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type DatabaseORM struct {
	PgxDB *config.PostgreDB
	RdsDB *config.RedisDB
}

func NewPgxConn(cfg *config.Config) *config.PostgreDB {
	dsn := "postgres://" + cfg.Database.PgxUser + ":" + cfg.Database.PgxPassword + "@"
	dsn = dsn + cfg.Database.PgxAddress + ":" + cfg.Database.PgxPort + "/"
	dsn = dsn + cfg.Database.PgxDB

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Postgre connection err - %v", err)
		return nil
	}
	return &config.PostgreDB{
		DB: db,
	}
}
func NewDBUsage(cfg *config.Config) *DatabaseORM {
	return &DatabaseORM{
		PgxDB: NewPgxConn(cfg),
		RdsDB: &config.RedisDB{
			RedisConnection: redis.Pool{
				Dial: func() (redis.Conn, error) {
					return redis.DialURL(cfg.Database.RedisURL)
				},
				MaxIdle:     8,
				MaxActive:   0,
				IdleTimeout: 100,
			},
		},
	}
}

func (db *DatabaseORM) GetAuthUserData(login string) (*config.UserAuthData, error) {
	rows, err := db.PgxDB.DB.Query("SELECT login, password FROM public.auth WHERE login=$1;", login)
	defer rows.Close()

	if err != nil {
		log.Printf("AuthorizeUser err - %v", err)
		return &config.UserAuthData{}, err
	}

	authData := &config.UserAuthData{}

	for rows.Next() {
		err = rows.Scan(&authData.Login, &authData.Password)
		if err != nil {
			log.Printf("PostgreORM.LoginUser scan error - %v", err)
			return &config.UserAuthData{}, err
		}
	}
	return authData, nil
}

func (db *DatabaseORM) RegisterUser(login, plainPassword string) (int, error) {
	rows, err := db.PgxDB.DB.Query("SELECT login, password FROM public.auth WHERE login=$1;", login)
	defer rows.Close()

	authData := &config.UserAuthData{}

	for rows.Next() {
		err = rows.Scan(&authData.Login, &authData.Password)
		if err == nil {
			log.Printf("User already exists - %v", err)
			return 0, ErrorUserAlreadyExists
		}
	}

	hashedPassword, err := argonpass.Hash(plainPassword, nil)
	if err != nil {
		log.Printf("Agron hashing err - %v\n", err)
		return 0, err
	}

	currentUsr := &config.UserAuthData{
		Login:    login,
		Password: hashedPassword,
	}

	_, err = db.PgxDB.DB.Exec("INSERT INTO public.auth(login, password) VALUES ($1, $2);", currentUsr.Login, currentUsr.Password)
	if err != nil {
		return 0, err
	}
	_, err = db.PgxDB.DB.Exec("INSERT INTO public.users(login) VALUES ($1);", currentUsr.Login)
	if err != nil {
		return 0, errors.New("error while insert new user to public.users")
	}

	var lastId int
	rows, _ = db.PgxDB.DB.Query("SELECT id FROM public.auth WHERE login=$1", currentUsr.Login)
	for rows.Next() {
		err = rows.Scan(&lastId)
		if err != nil {
			return 0, errors.New("Error while creating new user")
		}
	}

	return lastId, nil
}

func (db *DatabaseORM) CreateSession(in *config.Session) (*config.Session, error) {
	newSessionID := config.SessionID{ID: RandStringRunes(sessionIDKeyLen)}
	in.SessID = newSessionID
	dataSerialized, _ := json.Marshal(in)
	mKey := "SESSIONS: " + newSessionID.ID
	currentConn := db.RdsDB.RedisConnection.Get()
	result, err := redis.String(currentConn.Do("SET", mKey, dataSerialized, "EX", 259200))

	if err != nil {
		return nil, ErrorNewSession
	}
	if result != "OK" {
		return nil, ErrorResultNotOK
	}
	return in, nil
}

func (db *DatabaseORM) CheckSession(in *config.SessionID) (*config.Session, error) {
	mKey := "SESSIONS: " + in.ID
	currentConn := db.RdsDB.RedisConnection.Get()
	data, err := redis.Bytes(currentConn.Do("GET", mKey))
	if err != nil {
		log.Printf("redisORM error - %v", err)
		return &config.Session{}, err
	}
	currentSession := &config.Session{}
	err = json.Unmarshal(data, currentSession)
	if err != nil {
		log.Printf("Can't unpack session data to session interface")
		return &config.Session{}, err
	}
	return currentSession, nil
}

func (db *DatabaseORM) DeleteSession(in *config.SessionID) error {
	mKey := "SESSIONS:" + in.ID
	currentConn := db.RdsDB.RedisConnection.Get()
	_, err := redis.Int(currentConn.Do("DEL", mKey))
	if err != nil {
		log.Printf("Redis err | Deleting session err - %v:", err)
		return err
	}
	return nil
}

// Profile interactions

func (db *DatabaseORM) CreateProfile(in *config.User, username string) error {
	rows, err := db.PgxDB.DB.Query("SELECT first_name, second_name, sex, age, address, register_date, edit_date, login FROM public.users WHERE login=$1", username)
	defer rows.Close()
	if err != nil {
		return err
	}

	currentProfile := &config.User{}
	for rows.Next() {
		err = rows.Scan(&currentProfile.FirstName, &currentProfile.SecondName, &currentProfile.Sex, &currentProfile.Age,
			&currentProfile.Address, &currentProfile.RegisterDate, &currentProfile.EditDate, &currentProfile.Login)

		if err == nil {
			return errors.New("current user exists")
		}
	}

	_, err = db.PgxDB.DB.Exec("UPDATE public.users SET first_name=$2, second_name=$3, sex=$4, age=$5, address=$6, register_date=$7, edit_date=$8 WHERE login=$1;",
		username, in.FirstName, in.SecondName, in.Sex, in.Age,
		in.Address, in.RegisterDate, in.EditDate)
	if err != nil {
		log.Printf("Creating profile for user %v err - %v", username, err)
		return err
	}
	return nil
}

func (db *DatabaseORM) GetProfile(username string) (config.User, error) {
	rows, err := db.PgxDB.DB.Query("SELECT first_name, second_name, sex, age, address, register_date, edit_date, login FROM public.users WHERE login=$1", username)
	defer rows.Close()

	if err != nil {
		log.Printf("GetProfile error - %v", err)
		return config.User{}, err
	}

	currentProfile := &config.User{}

	for rows.Next() {
		err = rows.Scan(&currentProfile.FirstName, &currentProfile.SecondName, &currentProfile.Sex, &currentProfile.Age,
			&currentProfile.Address, &currentProfile.RegisterDate, &currentProfile.EditDate, &currentProfile.Login)

		if err != nil {
			log.Printf("Reading row in GetProfile err - %v", err)
			return config.User{}, err
		}
	}
	if currentProfile.FirstName == "" {
		log.Printf("There's no profile to user with this login")
		return config.User{}, errors.New("empty profile")
	}

	return *currentProfile, nil
}

func (db *DatabaseORM) EditProfile(username string, newData *config.User) error {
	_, err := db.PgxDB.DB.Exec("UPDATE public.users SET first_name=$2, second_name=$3, sex=$4, age=$5, address=$6, register_date=$7, edit_date=$8 WHERE login=$1;",
		username, newData.FirstName, newData.SecondName, newData.Sex, newData.Age, newData.Address, newData.RegisterDate, newData.EditDate)

	if err != nil {
		log.Printf("Edit profile for user %v err - %v", username, err)
		return err
	}
	return nil
}

func (db *DatabaseORM) DeleteProfile(username string) error {
	_, err := db.PgxDB.DB.Exec("UPDATE public.users SET first_name=NULL, second_name=NULL, sex=NULL, age=NULL, address=NULL, register_date=NULL, edit_date=NULL WHERE login=$1;", username)
	if err != nil {
		return err
	}
	return nil
}
