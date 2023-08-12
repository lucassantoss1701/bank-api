package mock

import (
	"context"
	"lucassantoss1701/bank/internal/usecase"

	"github.com/stretchr/testify/mock"
)

type FindTransfersByAccountUseCaseMock struct {
	mock.Mock
}

func NewFindTransfersByAccountUseCaseMock() *FindTransfersByAccountUseCaseMock {
	return &FindTransfersByAccountUseCaseMock{}
}

func (f *FindTransfersByAccountUseCaseMock) Execute(ctx context.Context, input *usecase.FindTransfersByAccountUseCaseInput) ([]usecase.FindTransfersByAccountUseCaseOutput, error) {
	args := f.Called(ctx, input)
	return args.Get(0).([]usecase.FindTransfersByAccountUseCaseOutput), args.Error(1)
}
