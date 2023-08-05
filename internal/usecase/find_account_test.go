package usecase_test

import (
	"errors"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/entity/mock"
	"lucassantoss1701/bank/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAccountUseCase_Execute(t *testing.T) {
	t.Run("Testing FindAccountUseCase when have success on find accounts", func(t *testing.T) {
		limit := 20
		offset := 0

		repository := mock.NewAccountRepositoryMock()
		accounts := mock.GetAccounts()
		repository.On("Find", limit, offset).Return(accounts, nil)

		findAccountUseCase := usecase.NewFindAccountUseCase(repository)

		input := usecase.NewFindAccountUseCaseInput(limit, offset)

		output, err := findAccountUseCase.Execute(input)

		assert.Nil(t, err)

		assert.NotEmpty(t, output)

		assert.Equal(t, accounts[0].ID, output[0].ID)
		assert.Equal(t, accounts[0].Name, output[0].Name)
		assert.Equal(t, accounts[0].CreatedAt, output[0].CreatedAt)

	})

	t.Run("Testing FindAccountUseCase when repository return a error", func(t *testing.T) {
		limit := 20
		offset := 0

		repository := mock.NewAccountRepositoryMock()
		repository.On("Find", limit, offset).Return([]entity.Account{}, errors.New("error on find orders"))

		findAccountUseCase := usecase.NewFindAccountUseCase(repository)

		input := usecase.NewFindAccountUseCaseInput(limit, offset)

		output, err := findAccountUseCase.Execute(input)

		assert.Empty(t, output)
		assert.NotNil(t, err)
		assert.Equal(t, "error on find orders", err.Error())

	})
}
