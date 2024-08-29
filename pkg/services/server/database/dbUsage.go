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
	redisAddress           = "redis://user:@localhost:6379/0"
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

// todo: в будущем перенести адрес редиса в .env файл
func NewRedisConn() *config.RedisDB {
	redisConn, err := redis.DialURL(redisAddress)
	if err != nil {
		log.Fatalf("Redis connection err - %v", err)
		return nil
	}
	return &config.RedisDB{
		RedisConnection: redisConn,
	}
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
		RdsDB: NewRedisConn(),
	}
}

func (db *DatabaseORM) GetUserByID(id int) (config.User, error) {
	item := &config.User{}
	rows, err := db.PgxDB.DB.Query(
		"SELECT user_id, first_name, second_name, sex, age, address, register_date, edit_date, login FROM public.users WHERE user_id=$1;", id)
	if err != nil {
		log.Printf("PostgreORM.GetUserById error - %v", err)
		return config.User{}, err
	}
	for rows.Next() {
		err = rows.Scan(&item.UserId, &item.FirstName, &item.SecondName, &item.Sex, &item.Age, &item.Address, &item.RegisterDate, &item.EditDate, &item.Login)
		if err != nil {
			return config.User{}, err
		}
	}
	return *item, nil
}

func (db *DatabaseORM) GetUserByLogin(login string) (config.User, error) {
	item := &config.User{}
	rows, err := db.PgxDB.DB.Query(
		"SELECT user_id, first_name, second_name, sex, age, address, register_date, edit_date, login FROM public.users WHERE login=$1;", login)
	if err != nil {
		log.Printf("PostgreORM.GetUserByLogin error - %v", err)
		return config.User{}, err
	}
	for rows.Next() {
		err = rows.Scan(&item.UserId, &item.FirstName, &item.SecondName, &item.Sex, &item.Age, &item.Address, &item.RegisterDate, &item.EditDate, &item.Login)
		if err != nil {
			return config.User{}, err
		}
	}
	return *item, nil
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

	currentUsr := &config.UserAuthData{
		Login:    login,
		Password: hashedPassword,
	}

	_, err = db.PgxDB.DB.Exec("INSERT INTO public.auth(login, password) VALUES ($1, $2);", currentUsr.Login, currentUsr.Password)
	if err != nil {
		return 0, err
	}
	_, err = db.PgxDB.DB.Exec("INSERT INTO public.users(login) VALUES ($1);", currentUsr.Login)

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
	result, err := redis.String(db.RdsDB.RedisConnection.Do("SET", mKey, dataSerialized, "EX", 259200))

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
	data, err := redis.Bytes(db.RdsDB.RedisConnection.Do("GET", mKey))
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
	_, err := redis.Int(db.RdsDB.RedisConnection.Do("DEL", mKey))
	if err != nil {
		log.Printf("Redis err | Deleting session err - %v:", err)
		return err
	}
	return nil
}
