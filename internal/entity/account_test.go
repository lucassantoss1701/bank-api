package entity_test

import (
	"lucassantoss1701/bank/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAccount_NewAccount(t *testing.T) {
	t.Run("Testing NewAccount when  returning a valid account", func(t *testing.T) {
		ID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
		name := "lucas"
		CPF := "35768297090"
		secret := "4578405"
		balance := 100
		createdAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
		account, err := entity.NewAccount(ID, name, CPF, secret, balance, &createdAt)

		assert.Nil(t, err)
		assert.NotNil(t, account)

		assert.Equal(t, ID, account.ID)
		assert.Equal(t, name, account.Name)
		assert.Equal(t, CPF, account.CPF)
		assert.Equal(t, balance, account.Balance)
		assert.Equal(t, &createdAt, account.CreatedAt)

		assert.NotEqual(t, secret, account.Secret)
		assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(account.Secret), []byte(secret)))

	})

	t.Run("Testing NewAccount when returning an invalid account (too loong secret)", func(t *testing.T) {
		ID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
		name := "lucas"
		CPF := "35768297090"
		secret := "45784054578405457840545784054578405454578405457840545784054578405457840578405457840545784054578405457840545784054578405"
		balance := 100
		createdAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
		account, err := entity.NewAccount(ID, name, CPF, secret, balance, &createdAt)

		assert.Nil(t, account)
		assert.NotNil(t, err)

		assert.Equal(t, "error on hashing password", err.Error())
	})

	t.Run("Testing NewAccount when returning an invalid account (ID is invalid)", func(t *testing.T) {
		ID := ""
		name := "lucas"
		CPF := "35768297090"
		secret := "4578405"
		balance := 100
		createdAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
		account, err := entity.NewAccount(ID, name, CPF, secret, balance, &createdAt)

		assert.Nil(t, account)
		assert.NotNil(t, err)

		assert.Equal(t, "ID cannot be empty", err.Error())
	})

	t.Run("Testing NewAccount when returning an invalid account (Name is invalid)", func(t *testing.T) {
		ID := "123"
		name := ""
		CPF := "35768297090"
		secret := "4578405"
		balance := 100
		createdAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
		account, err := entity.NewAccount(ID, name, CPF, secret, balance, &createdAt)

		assert.Nil(t, account)
		assert.NotNil(t, err)

		assert.Equal(t, "name cannot be empty", err.Error())
	})

	t.Run("Testing NewAccount when returning an invalid account (CPF is invalid)", func(t *testing.T) {
		ID := "123"
		name := "Lucas"
		CPF := ""
		secret := "4578405"
		balance := 100
		createdAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
		account, err := entity.NewAccount(ID, name, CPF, secret, balance, &createdAt)

		assert.Nil(t, account)
		assert.NotNil(t, err)

		assert.Equal(t, "CPF cannot be empty", err.Error())
	})

	t.Run("Testing NewAccount when returning an invalid account (Secret is invalid)", func(t *testing.T) {
		ID := "123"
		name := "Lucas"
		CPF := "35768297090"
		secret := ""
		balance := 100
		createdAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
		account, err := entity.NewAccount(ID, name, CPF, secret, balance, &createdAt)

		assert.Nil(t, account)
		assert.NotNil(t, err)

		assert.Equal(t, "secret cannot be empty", err.Error())
	})

	t.Run("Testing NewAccount when returning an invalid account (Balance is invalid)", func(t *testing.T) {
		ID := "123"
		name := "Lucas"
		CPF := "35768297090"
		secret := "4578405"
		balance := -100
		createdAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
		account, err := entity.NewAccount(ID, name, CPF, secret, balance, &createdAt)

		assert.Nil(t, account)
		assert.NotNil(t, err)

		assert.Equal(t, "balance cannot be minor than 0", err.Error())
	})

	t.Run("Testing NewAccount when returning an invalid account (CreatedAt is invalid)", func(t *testing.T) {
		ID := "123"
		name := "Lucas"
		CPF := "35768297090"
		secret := "4578405"
		balance := 100
		var createdAt *time.Time
		account, err := entity.NewAccount(ID, name, CPF, secret, balance, createdAt)

		assert.Nil(t, account)
		assert.NotNil(t, err)

		assert.Equal(t, "created at cannot be nil", err.Error())
	})

}
