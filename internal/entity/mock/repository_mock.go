package mock

import (
	"context"
	"lucassantoss1701/bank/internal/entity"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{}
}

func (r *RepositoryMock) BeginTx(ctx context.Context) (entity.TransactionHandler, error) {
	args := r.Called(ctx)

	return args.Get(0).(entity.TransactionHandler), args.Error(1)
}

func (r *RepositoryMock) CommitTx(tx entity.TransactionHandler) error {
	args := r.Called(tx)
	return args.Error(0)

}
func (r *RepositoryMock) RollbackTx(tx entity.TransactionHandler) error {
	args := r.Called(tx)
	return args.Error(0)
}
