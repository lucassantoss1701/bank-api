package entity

import (
	"context"
	"database/sql"
)

type AccountRepository interface {
	Find(ctx context.Context, limit, offset int) ([]Account, error)
	FindByID(ctx context.Context, ID string) (Account, error)
	Create(ctx context.Context, account *Account) (Account, error)
	UpdateBalance(ctx context.Context, accountID string, newBalance int, tx ...TransactionHandler) (Account, error)
}

type TransferRepository interface {
	FindByAccountID(ctx context.Context, AccountID string, limit, offset int) ([]Transfer, error)
	Create(ctx context.Context, transfer *Transfer, tx ...TransactionHandler) (Transfer, error)
}

type Repository interface {
	BeginTx(ctx context.Context) (TransactionHandler, error)
	CommitTx(tx TransactionHandler) error
	RollbackTx(tx TransactionHandler) error
}

type TransactionHandler interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
