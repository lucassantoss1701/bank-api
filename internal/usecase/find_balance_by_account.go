package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
)

type FindBalanceByAccountUseCase struct {
	repostiory entity.AccountRepository
}

func NewFindBalanceByAccountUseCase(repostiory entity.AccountRepository) *FindBalanceByAccountUseCase {
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
