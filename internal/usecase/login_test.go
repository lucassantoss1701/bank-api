package usecase_test

import (
	"context"
	"errors"
	"lucassantoss1701/bank/internal/entity/mock"
	"lucassantoss1701/bank/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUseCase_Execute(t *testing.T) {
	t.Run("Testing LoginUseCase when have success on login", func(t *testing.T) {
		ctx := context.Background()
		repository := mock.NewAccountRepositoryMock()
		account := mock.CreateAccount()
		loginUseCase := usecase.NewLoginUseCase(repository)

		secretHashed, err := bcrypt.GenerateFromPassword([]byte(account.Secret), bcrypt.DefaultCost)
		assert.Nil(t, err)

		account.Secret = string(secretHashed)

		CPF := "34688151071"
		secret := "5e0542f964858f96ae7194fb2a7dd365"

		repository.On("FindByCPF", ctx, CPF).Return(account, nil)

		input := usecase.NewLoginUseCaseInput(CPF, secret, "")

		output, err := loginUseCase.Execute(ctx, input)
		assert.Nil(t, err)

		assert.NotEmpty(t, output)
		assert.NotEmpty(t, output.Token)
	})

	t.Run("Testing LoginUseCase when repository returns an error", func(t *testing.T) {
		ctx := context.Background()
		repository := mock.NewAccountRepositoryMock()
		account := mock.CreateAccount()
		loginUseCase := usecase.NewLoginUseCase(repository)

		secretHashed, err := bcrypt.GenerateFromPassword([]byte(account.Secret), bcrypt.DefaultCost)
		assert.Nil(t, err)

		account.Secret = string(secretHashed)

		CPF := "34688151071"
		secret := "5e0542f964858f96ae7194fb2a7dd365"

		repository.On("FindByCPF", ctx, CPF).Return(account, errors.New("error on find account"))

		input := usecase.NewLoginUseCaseInput(CPF, secret, "")

		_, err = loginUseCase.Execute(ctx, input)
		assert.NotNil(t, err)

		assert.Equal(t, "error on find account", err.Error())
	})

	t.Run("Testing LoginUseCase when secret is incorret", func(t *testing.T) {
		ctx := context.Background()
		repository := mock.NewAccountRepositoryMock()
		account := mock.CreateAccount()
		loginUseCase := usecase.NewLoginUseCase(repository)

		secretHashed, err := bcrypt.GenerateFromPassword([]byte(account.Secret), bcrypt.DefaultCost)
		assert.Nil(t, err)

		account.Secret = string(secretHashed)

		CPF := "34688151071"
		secret := "incorret secret"

		repository.On("FindByCPF", ctx, CPF).Return(account, nil)

		input := usecase.NewLoginUseCaseInput(CPF, secret, "")

		_, err = loginUseCase.Execute(ctx, input)
		assert.NotNil(t, err)

		assert.Equal(t, "secret is incorrect", err.Error())
	})
}
