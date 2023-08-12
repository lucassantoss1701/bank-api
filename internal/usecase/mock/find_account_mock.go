package mock

import (
	"context"
	"lucassantoss1701/bank/internal/usecase"

	"github.com/stretchr/testify/mock"
)

type FindAccountUseCaseMock struct {
	mock.Mock
}

func NewFindAccountUseCaseMock() *FindAccountUseCaseMock {
	return &FindAccountUseCaseMock{}
}

func (f *FindAccountUseCaseMock) Execute(ctx context.Context, input *usecase.FindAccountUseCaseInput) ([]usecase.FindAccountUseCaseOutput, error) {
	args := f.Called(ctx, input)
	return args.Get(0).([]usecase.FindAccountUseCaseOutput), args.Error(1)
}
