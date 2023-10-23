package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

// HashPassword return the bcrypt hash of the password
func HashedPassword(password string) (string, error) {
	byteHashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(byteHashedPassword), nil
}
