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

func TestFindTransferencesByAccountUseCase_Execute(t *testing.T) {
	t.Run("Testing FindTransferencesByAccountUseCase when have success on find balance", func(t *testing.T) {
		ctx := context.Background()
		accountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
		limit := 20
		offset := 0

		repository := mock.NewTransferRepositoryMock()
		transferences := mock.GetTransfererences()
		repository.On("FindByAccountID", ctx, accountID, limit, offset).Return(transferences, nil)

		findTransferencesByAccountUseCase := usecase.NewFindTransferencesByAccountUseCase(repository)

		input := usecase.NewFindTransferencesByAccountUseCaseInput(accountID, limit, offset)

		output, err := findTransferencesByAccountUseCase.Execute(ctx, input)

		assert.Nil(t, err)
		assert.NotNil(t, output)

		assert.Equal(t, transferences[0].ID, output[0].ID)
		assert.Equal(t, transferences[0].Amount, output[0].Amount)
		assert.Equal(t, transferences[0].CreatedAt, output[0].CreatedAt)
		assert.Equal(t, transferences[0].DestinationAccount.Name, output[0].DestinationAccount.Name)
		assert.Equal(t, transferences[0].DestinationAccount.ID, output[0].DestinationAccount.ID)
		assert.Equal(t, transferences[0].DestinationAccount.CreatedAt, output[0].DestinationAccount.CreatedAt)

	})

	t.Run("Testing FindTransferencesByAccountUseCase when repository returns an error", func(t *testing.T) {
		ctx := context.Background()

		accountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
		limit := 20
		offset := 0

		repository := mock.NewTransferRepositoryMock()
		repository.On("FindByAccountID", ctx, accountID, limit, offset).Return([]entity.Transfer{}, errors.New("error on find transferecenes by account ID"))

		findTransferencesByAccountUseCase := usecase.NewFindTransferencesByAccountUseCase(repository)

		input := usecase.NewFindTransferencesByAccountUseCaseInput(accountID, limit, offset)

		output, err := findTransferencesByAccountUseCase.Execute(ctx, input)

		assert.Nil(t, output)

		assert.NotNil(t, err)
		assert.Equal(t, "error on find transferecenes by account ID", err.Error())

	})

}
