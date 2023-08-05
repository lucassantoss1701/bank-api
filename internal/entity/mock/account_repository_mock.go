package mock

import (
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

func (a *AccountRepositoryMock) Find(limit, offset int) ([]entity.Account, error) {
	args := a.Called(limit, offset)
	return args.Get(0).([]entity.Account), args.Error(1)

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
