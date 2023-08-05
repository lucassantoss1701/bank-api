package usecase

import "lucassantoss1701/bank/internal/entity"

type FindBalanceByAccountUseCase struct {
	repostiory entity.AccountRepositoryInterface
}

func NewFindBalanceByAccountUseCase(repostiory entity.AccountRepositoryInterface) *FindBalanceByAccountUseCase {
	return &FindBalanceByAccountUseCase{
		repostiory: repostiory,
	}
}

func (g *FindBalanceByAccountUseCase) Execute(ID string) (*FindBalanceByAccountUseCaseOutput, error) {
	account, err := g.repostiory.FindByID(ID)
	if err != nil {
		return nil, err
	}
	return &FindBalanceByAccountUseCaseOutput{Balance: account.Balance}, nil

}

type FindBalanceByAccountUseCaseOutput struct {
	Balance int
}
