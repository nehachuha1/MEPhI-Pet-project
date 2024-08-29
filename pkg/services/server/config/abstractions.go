package config

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
)

type User struct {
	UserId       int    `json:"user_id"`
	Login        string `json:"login"`
	FirstName    string `json:"first_name"`
	SecondName   string `json:"second_name"`
	Sex          string `json:"sex"`
	Age          int    `json:"age"`
	Address      string `json:"address"`
	RegisterDate string `json:"registerDate"`
	EditDate     string `json:"editDate"`
}

type UserAuthData struct {
	Login    string
	Password string
}

type Session struct {
	SessID   SessionID
	Username string
}

type SessionID struct {
	ID string
}

type PostgreDB struct {
	DB *sql.DB
}

type RedisDB struct {
	RedisConnection redis.Conn
}
