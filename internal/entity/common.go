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
