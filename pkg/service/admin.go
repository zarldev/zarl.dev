package service

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func NewHashedPassword(pass string) (string, error) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return pass, err
	}
	return string(bcryptPassword), nil
}

func RandomUsernameAndHashedPassword() (string, string) {
	username := randomString(16)
	password := randomString(16)
	password, _ = NewHashedPassword(password)
	return username, password
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
