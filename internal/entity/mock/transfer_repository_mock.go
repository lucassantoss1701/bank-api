package mock

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"

	"github.com/stretchr/testify/mock"
)

type TransfersRepositoryMock struct {
	mock.Mock
}

func NewTransferRepositoryMock() *TransfersRepositoryMock {
	return &TransfersRepositoryMock{}
}

func (a *TransfersRepositoryMock) FindByAccountID(ctx context.Context, accountID string, limit, offset int) ([]entity.Transfer, error) {
	args := a.Called(ctx, accountID, limit, offset)
	return args.Get(0).([]entity.Transfer), args.Error(1)
}

func (t *TransfersRepositoryMock) Create(ctx context.Context, transfer *entity.Transfer, tx ...entity.TransactionHandler) (entity.Transfer, error) {
	args := t.Called(ctx, transfer, tx)
	return args.Get(0).(entity.Transfer), args.Error(1)
}

func GetTransfererences() []entity.Transfer {
	date := time.Date(2023, 8, 5, 16, 00, 00, 00, time.UTC)

	originAccountDate := time.Date(2023, 2, 5, 13, 00, 00, 00, time.UTC)
	destinationAccountDate := time.Date(2023, 1, 4, 12, 00, 00, 00, time.UTC)

	return []entity.Transfer{
		{
			ID:        "2bd765a6-47bd-4731-9eb2-1e65542f4477",
			Amount:    500,
			CreatedAt: &date,
			OriginAccount: &entity.Account{
				ID:        "2bd765a6-47bd-4731-9eb2-1e65542f4477",
				Name:      "Lucas",
				CPF:       "",
				Secret:    "",
				Balance:   0,
				CreatedAt: &originAccountDate,
			},
			DestinationAccount: &entity.Account{
				ID:        "6ac7ebbf-568b-45f2-a295-bfbab73f1cf6",
				Name:      "Rogerio",
				CPF:       "",
				Secret:    "",
				Balance:   0,
				CreatedAt: &destinationAccountDate,
			},
		},
	}
}
