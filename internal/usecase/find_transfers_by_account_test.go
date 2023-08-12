package usecase_test

import (
	"context"
	"errors"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/entity/mock"
	"lucassantoss1701/bank/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindTransfersByAccountUseCase_Execute(t *testing.T) {
	t.Run("Testing FindTransfersByAccountUseCase when have success on find balance", func(t *testing.T) {
		ctx := context.Background()
		accountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
		limit := 20
		offset := 0

		repository := mock.NewTransferRepositoryMock()
		transfers := mock.GetTransfererences()
		repository.On("FindByAccountID", ctx, accountID, limit, offset).Return(transfers, nil)

		findTransfersByAccountUseCase := usecase.NewFindTransfersByAccountUseCase(repository)

		input := usecase.NewFindTransfersByAccountUseCaseInput(accountID, limit, offset)

		output, err := findTransfersByAccountUseCase.Execute(ctx, input)

		assert.Nil(t, err)
		assert.NotNil(t, output)

		assert.Equal(t, transfers[0].ID, output[0].ID)
		assert.Equal(t, transfers[0].Amount, output[0].Amount)
		assert.Equal(t, transfers[0].CreatedAt.Format(time.RFC3339), output[0].CreatedAt)
		assert.Equal(t, transfers[0].DestinationAccount.Name, output[0].DestinationAccount.Name)
		assert.Equal(t, transfers[0].DestinationAccount.ID, output[0].DestinationAccount.ID)

	})

	t.Run("Testing FindTransfersByAccountUseCase when repository returns an error", func(t *testing.T) {
		ctx := context.Background()

		accountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
		limit := 20
		offset := 0

		repository := mock.NewTransferRepositoryMock()
		repository.On("FindByAccountID", ctx, accountID, limit, offset).Return([]entity.Transfer{}, errors.New("error on find transfer by account ID"))

		findTransfersByAccountUseCase := usecase.NewFindTransfersByAccountUseCase(repository)

		input := usecase.NewFindTransfersByAccountUseCaseInput(accountID, limit, offset)

		output, err := findTransfersByAccountUseCase.Execute(ctx, input)

		assert.Nil(t, output)

		assert.NotNil(t, err)
		assert.Equal(t, "error on find transfer by account ID", err.Error())

	})

}
