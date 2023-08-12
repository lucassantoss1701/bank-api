package mock

import (
	"context"
	"lucassantoss1701/bank/internal/usecase"

	"github.com/stretchr/testify/mock"
)

type FindBalanceByAccountUseCaseMock struct {
	mock.Mock
}

func NewFindBalanceByAccountUseCaseMock() *FindBalanceByAccountUseCaseMock {
	return &FindBalanceByAccountUseCaseMock{}
}

func (c *FindBalanceByAccountUseCaseMock) Execute(ctx context.Context, input *usecase.FindBalanceByAccountUseCaseInput) (*usecase.FindBalanceByAccountUseCaseOutput, error) {
	args := c.Called(ctx, input)
	return args.Get(0).(*usecase.FindBalanceByAccountUseCaseOutput), args.Error(1)
}
