package config

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"html/template"
)

type User struct {
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
	RedisConnection redis.Pool
}

type Templates struct {
	Templates *template.Template
}
