package mock

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type TransactionHandlerMock struct {
	mock.Mock
}

func NewTransactionHandlerMock() *TransactionHandlerMock {
	return &TransactionHandlerMock{}
}

func (t *TransactionHandlerMock) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	argsMock := t.Called(ctx, query, args)
	return argsMock.Get(0).(sql.Result), argsMock.Error(1)
}

func (t *TransactionHandlerMock) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	argsMock := t.Called(ctx, query, args)
	return argsMock.Get(0).(*sql.Row)
}
