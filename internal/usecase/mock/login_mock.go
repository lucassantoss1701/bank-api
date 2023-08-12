package mock

import (
	"context"
	"lucassantoss1701/bank/internal/usecase"

	"github.com/stretchr/testify/mock"
)

type LoginUseCaseMock struct {
	mock.Mock
}

func NewLoginUseCaseMock() *LoginUseCaseMock {
	return &LoginUseCaseMock{}
}

func (f *LoginUseCaseMock) Execute(ctx context.Context, input *usecase.LoginUseCaseInput) (*usecase.LoginUseCaseOutput, error) {
	args := f.Called(ctx, input)
	return args.Get(0).(*usecase.LoginUseCaseOutput), args.Error(1)
}
