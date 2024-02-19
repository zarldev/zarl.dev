package service

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func NewHashedPassword(pass string) (string, error) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return pass, err
	}
	return string(bcryptPassword), nil
}

func AdminPassCrypted() (string, string, string) {
	username := randomString(16)
	password := randomString(16)
	hashedPassword, _ := NewHashedPassword(password)
	return username, password, hashedPassword
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		max := big.NewInt(int64(len(letters)))
		randIndex, err := rand.Int(rand.Reader, max)
		if err != nil {
			return ""
		}
		b[i] = letters[randIndex.Int64()]
	}
	return string(b)
}
