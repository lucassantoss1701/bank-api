package entity

type AccountRepositoryInterface interface {
	Find(limit, offset int) ([]Account, error)
	FindByID(ID string) (Account, error)
	Create(account *Account) (Account, error)
}

type TransferRepositoryInterface interface {
	FindByAccountID(AccountID string, limit, offset int) ([]Transfer, error)
}
