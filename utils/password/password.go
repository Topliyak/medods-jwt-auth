package password

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/medods-jwt-auth/config"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), config.BCRYPT_COAST)
	return string(hash), err
}

func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
