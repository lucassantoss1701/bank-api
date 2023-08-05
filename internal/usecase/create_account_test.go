package usecase_test

import (
	"errors"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/entity/mock"
	"lucassantoss1701/bank/internal/usecase"
	"testing"

	testify "github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccountUseCase_Execute(t *testing.T) {
	t.Run("Testing CreateAccountUseCase when have success on create account", func(t *testing.T) {

		repository := mock.NewAccountRepositoryMock()
		account := mock.CreateAccount()
		repository.On("Create", testify.AnythingOfTypeArgument("*entity.Account")).Return(account, nil)

		createAccountUseCase := usecase.NewCreateAccountUseCase(repository)

		input := usecase.NewCreateAccountUseCaseInput(account.ID, account.Name, account.CPF, account.Secret, account.Balance, *account.CreatedAt)

		output, err := createAccountUseCase.Execute(input)

		assert.Nil(t, err)

		assert.NotEmpty(t, output)

		assert.Equal(t, account.ID, output.ID)
		assert.Equal(t, account.Name, output.Name)
		assert.Equal(t, account.Balance, output.Balance)
		assert.Equal(t, account.CreatedAt, output.CreatedAt)

	})

	t.Run("Testing CreateAccountUseCase when create new account return error", func(t *testing.T) {

		repository := mock.NewAccountRepositoryMock()
		account := mock.CreateAccount()
		repository.On("Create", testify.AnythingOfTypeArgument("*entity.Account")).Return(account, nil)

		createAccountUseCase := usecase.NewCreateAccountUseCase(repository)

		input := usecase.NewCreateAccountUseCaseInput("", account.Name, account.CPF, account.Secret, account.Balance, *account.CreatedAt)

		output, err := createAccountUseCase.Execute(input)

		assert.Nil(t, output)

		assert.NotNil(t, err)
		assert.Equal(t, "ID cannot be empty", err.Error())

	})

	t.Run("Testing CreateAccountUseCase when repository return error", func(t *testing.T) {

		repository := mock.NewAccountRepositoryMock()
		account := mock.CreateAccount()
		repository.On("Create", testify.AnythingOfTypeArgument("*entity.Account")).Return(entity.Account{}, errors.New("error on create account"))

		createAccountUseCase := usecase.NewCreateAccountUseCase(repository)

		input := usecase.NewCreateAccountUseCaseInput(account.ID, account.Name, account.CPF, account.Secret, account.Balance, *account.CreatedAt)

		output, err := createAccountUseCase.Execute(input)

		assert.Nil(t, output)

		assert.NotNil(t, err)
		assert.Equal(t, "error on create account", err.Error())

	})
}
