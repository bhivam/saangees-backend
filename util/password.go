package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassowrd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	return string(hashedPassowrd), nil
}

func CheckPassword(password string, hashedPassword string) error {
  err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
  if err != nil {
    return fmt.Errorf("invalid password: %v", err)
  }
  return nil
}
