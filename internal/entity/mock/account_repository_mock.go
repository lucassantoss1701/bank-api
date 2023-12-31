package mock

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"

	"github.com/stretchr/testify/mock"
)

type AccountRepositoryMock struct {
	mock.Mock
}

func NewAccountRepositoryMock() *AccountRepositoryMock {
	return &AccountRepositoryMock{}
}

func (a *AccountRepositoryMock) Find(ctx context.Context, limit, offset int) ([]entity.Account, error) {
	args := a.Called(ctx, limit, offset)
	return args.Get(0).([]entity.Account), args.Error(1)
}

func (a *AccountRepositoryMock) FindByID(ctx context.Context, ID string) (entity.Account, error) {
	args := a.Called(ctx, ID)
	return args.Get(0).(entity.Account), args.Error(1)
}

func (a *AccountRepositoryMock) Create(ctx context.Context, account *entity.Account) (entity.Account, error) {
	args := a.Called(ctx, account)
	return args.Get(0).(entity.Account), args.Error(1)
}

func (a *AccountRepositoryMock) UpdateBalance(ctx context.Context, accountID string, newBalance int, tx ...entity.TransactionHandler) (entity.Account, error) {
	args := a.Called(ctx, accountID, newBalance)
	return args.Get(0).(entity.Account), args.Error(1)
}

func (a *AccountRepositoryMock) FindByCPF(ctx context.Context, CPF string) (entity.Account, error) {
	args := a.Called(ctx, CPF)
	return args.Get(0).(entity.Account), args.Error(1)
}

func GetAccounts() []entity.Account {
	date := time.Date(2023, 8, 5, 16, 00, 00, 00, time.UTC)
	return []entity.Account{
		{
			ID:        "2bd765a6-47bd-4731-9eb2-1e65542f4477",
			Name:      "Lucas",
			CPF:       "",
			Secret:    "",
			Balance:   0,
			CreatedAt: &date,
		},
	}
}

func CreateAccount() entity.Account {
	date := time.Date(2023, 8, 5, 16, 00, 00, 00, time.UTC)
	return entity.Account{
		ID:        "2bd765a6-47bd-4731-9eb2-1e65542f4477",
		Name:      "Lucas",
		CPF:       "34688151071",
		Secret:    "5e0542f964858f96ae7194fb2a7dd365",
		Balance:   500,
		CreatedAt: &date,
	}
}
