package usecase_test

import (
	"context"
	"errors"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/entity/mock"
	"lucassantoss1701/bank/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindBalanceByAccountUseCase_Execute(t *testing.T) {
	t.Run("Testing FindBalanceByAccountUseCase when have success on find balance", func(t *testing.T) {
		ctx := context.Background()

		ID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"

		input := usecase.NewFindBalanceByAccountUseCaseInput(ID)

		repository := mock.NewAccountRepositoryMock()
		account := mock.GetAccounts()[0]
		repository.On("FindByID", ctx, ID).Return(account, nil)

		findBalanceByAccountUseCase := usecase.NewFindBalanceByAccountUseCase(repository)

		output, err := findBalanceByAccountUseCase.Execute(ctx, input)

		assert.Nil(t, err)

		assert.Equal(t, account.Balance, output.Balance)
	})

	t.Run("Testing FindBalanceByAccountUseCase when repository returns an error", func(t *testing.T) {
		ctx := context.Background()

		ID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
		input := usecase.NewFindBalanceByAccountUseCaseInput(ID)

		repository := mock.NewAccountRepositoryMock()
		repository.On("FindByID", ctx, ID).Return(entity.Account{}, errors.New("error on find account"))

		findBalanceByAccountUseCase := usecase.NewFindBalanceByAccountUseCase(repository)

		output, err := findBalanceByAccountUseCase.Execute(ctx, input)

		assert.NotNil(t, err)
		assert.Equal(t, "error on find account", err.Error())
		assert.Nil(t, output)
	})

}
