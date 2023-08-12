package entity

import (
	"errors"
	"time"
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

	if ID == "" {
		ID = NewUUID()
	}

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

func (a *Account) isValid() error {
	validationError := NewErrorHandler(ENTITY_ERROR)

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

func (a *Account) addFromBalance(value int) {
	a.Balance += value
}

func (a *Account) removeFromBalance(value int) error {
	validationError := NewErrorHandler(ENTITY_ERROR)

	a.Balance -= value

	if a.Balance < 0 {
		validationError.Add("new balance cannot be minor than 0")
		validationError.TypeError = ENTITY_ERROR
		return validationError
	}

	return nil
}
