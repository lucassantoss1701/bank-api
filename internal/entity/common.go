package entity

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUUID() string {
	return uuid.NewString()
}

func hash(value string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
}

func hashIsValid(hashedPassowrd string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassowrd), []byte(password))
	return err == nil
}
