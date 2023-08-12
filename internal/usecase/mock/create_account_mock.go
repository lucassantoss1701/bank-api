package mock

import (
	"context"
	"lucassantoss1701/bank/internal/usecase"

	"github.com/stretchr/testify/mock"
)

type CreateAccountUseCaseMock struct {
	mock.Mock
}

func NewCreateAccountUseCaseMock() *CreateAccountUseCaseMock {
	return &CreateAccountUseCaseMock{}
}

func (c *CreateAccountUseCaseMock) Execute(ctx context.Context, input *usecase.CreateAccountUseCaseInput) (*usecase.CreateAccountUseCaseOutput, error) {
	args := c.Called(ctx, input)
	return args.Get(0).(*usecase.CreateAccountUseCaseOutput), args.Error(1)
}
