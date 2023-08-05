package entity

type AccountRepositoryInterface interface {
	Find(limit, offset int) ([]Account, error)
}
