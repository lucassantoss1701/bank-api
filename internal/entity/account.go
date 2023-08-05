package entity

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID        string
	Name      string
	CPF       string
	Secret    string
	Balance   int
	CreatedAt *time.Time
}

func NewAccount(ID string, name string, CPF string, secret string, balance int, createdAt *time.Time) (*Account, error) {

	account := &Account{
		ID:        ID,
		Name:      name,
		CPF:       CPF,
		Secret:    secret,
		Balance:   balance,
		CreatedAt: createdAt,
	}

	err := account.isValid()
	if err != nil {
		return nil, err
	}

	secretHashed, err := hash(secret)
	if err != nil {
		return nil, errors.New("error on hashing password")
	}

	account.Secret = string(secretHashed)

	return account, nil

}

func hash(secret string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
}

func (a *Account) isValid() error {
	validationError := &ValidationError{}

	if a.ID == "" {
		validationError.Add("ID cannot be empty")
	}

	if a.Name == "" {
		validationError.Add("name cannot be empty")
	}

	if a.CPF == "" {
		validationError.Add("CPF cannot be empty")
	}

	if a.Secret == "" {
		validationError.Add("secret cannot be empty")
	}

	if a.Balance < 0 {
		validationError.Add("balance cannot be minor than 0")
	}

	if a.CreatedAt == nil {
		validationError.Add("created at cannot be nil")
	}

	if len(validationError.Messages) > 0 {
		return validationError
	}

	return nil
}

func (a *Account) hasSufficientBalance(value int) bool {
	return a.Balance-value > 0
}

func (a *Account) setBalance(value int) error {
	validationError := &ValidationError{}

	a.Balance += value

	if a.Balance < 0 {
		validationError.Add("new balance cannot be minor than 0")
		return validationError
	}

	return nil

}
