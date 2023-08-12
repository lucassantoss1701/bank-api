package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
)

type IFindBalanceByAccountUseCase interface {
	Execute(ctx context.Context, input *FindBalanceByAccountUseCaseInput) (*FindBalanceByAccountUseCaseOutput, error)
}

type FindBalanceByAccountUseCase struct {
	repostiory entity.AccountRepository
}

func NewFindBalanceByAccountUseCase(repostiory entity.AccountRepository) *FindBalanceByAccountUseCase {
	return &FindBalanceByAccountUseCase{
		repostiory: repostiory,
	}
}

func (g *FindBalanceByAccountUseCase) Execute(ctx context.Context, input *FindBalanceByAccountUseCaseInput) (*FindBalanceByAccountUseCaseOutput, error) {
	account, err := g.repostiory.FindByID(ctx, input.id)
	if err != nil {
		return nil, err
	}
	return NewFindBalanceByAccountUseCaseOutput(account.Balance), nil

}

type FindBalanceByAccountUseCaseInput struct {
	id string
}

func NewFindBalanceByAccountUseCaseInput(id string) *FindBalanceByAccountUseCaseInput {
	return &FindBalanceByAccountUseCaseInput{
		id: id,
	}
}

type FindBalanceByAccountUseCaseOutput struct {
	Balance int `json:"balance"`
}

func NewFindBalanceByAccountUseCaseOutput(balance int) *FindBalanceByAccountUseCaseOutput {
	return &FindBalanceByAccountUseCaseOutput{
		Balance: balance,
	}
}
