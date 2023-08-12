package mock

import (
	"context"
	"lucassantoss1701/bank/internal/usecase"

	"github.com/stretchr/testify/mock"
)

type MakeTransferUseCaseMock struct {
	mock.Mock
}

func NewMakeTransferUseCaseMock() *MakeTransferUseCaseMock {
	return &MakeTransferUseCaseMock{}
}

func (f *MakeTransferUseCaseMock) Execute(ctx context.Context, input *usecase.MakeTransferUseCaseInput) (*usecase.MakeTransferUseCaseOutput, error) {
	args := f.Called(ctx, input)
	return args.Get(0).(*usecase.MakeTransferUseCaseOutput), args.Error(1)
}
