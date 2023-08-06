package entity

import "context"

type AccountRepositoryInterface interface {
	Find(ctx context.Context, limit, offset int) ([]Account, error)
	FindByID(ctx context.Context, ID string) (Account, error)
	Create(ctx context.Context, account *Account) (Account, error)
}

type TransferRepositoryInterface interface {
	FindByAccountID(ctx context.Context, AccountID string, limit, offset int) ([]Transfer, error)
}
