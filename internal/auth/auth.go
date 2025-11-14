package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal("Error while creating password hash %v", err)
		return "", err
	}

	return hash, nil
}

func CheckPasswordHash(textPassword string, hash string) (bool, error) {
	valid, _, err := argon2id.CheckHash(textPassword, hash)
	if err != nil {
		log.Fatal("Error while Checking password hash %v", err)
		return false, err
	}

	return valid, nil
}
