package utils

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

// todo: Дописать авторизацию по примеру из лекций с алгоритмом Argon2

func EncryptPasswordWithSalt(plainPassword string, salt []byte) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func EncryptPassword(plainPassword string) (hashedPassword []byte) {
	salt := make([]byte, 8)
	rand.Read(salt)

	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func CheckPassword(dbPassword string, plainPassword string) bool {
	dbPasswordBytes := []byte(dbPassword)
	salt := dbPasswordBytes[0:8]
	hashedInputPassword := EncryptPasswordWithSalt(plainPassword, salt)
	return bytes.Equal(dbPasswordBytes, hashedInputPassword)
}
