package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
)

type FindBalanceByAccountUseCase struct {
	repostiory entity.AccountRepositoryInterface
}

func NewFindBalanceByAccountUseCase(repostiory entity.AccountRepositoryInterface) *FindBalanceByAccountUseCase {
	return &FindBalanceByAccountUseCase{
		repostiory: repostiory,
	}
}

func (g *FindBalanceByAccountUseCase) Execute(ctx context.Context, ID string) (*FindBalanceByAccountUseCaseOutput, error) {
	account, err := g.repostiory.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &FindBalanceByAccountUseCaseOutput{Balance: account.Balance}, nil

}

type FindBalanceByAccountUseCaseOutput struct {
	Balance int
}
