package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"
)

type CreateAccountUseCase struct {
	repostiory entity.AccountRepositoryInterface
}

func NewCreateAccountUseCase(repostiory entity.AccountRepositoryInterface) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		repostiory: repostiory,
	}
}

func (c *CreateAccountUseCase) Execute(ctx context.Context, input *CreateAccountUseCaseInput) (*CreateAccountUseCaseOutput, error) {
	account, err := entity.NewAccount(input.id, input.name, input.cpf, input.secret, input.balance, input.createdAt)
	if err != nil {
		return nil, err
	}

	createdAccount, err := c.repostiory.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return &CreateAccountUseCaseOutput{
		ID:        createdAccount.ID,
		Name:      createdAccount.Name,
		Balance:   createdAccount.Balance,
		CreatedAt: createdAccount.CreatedAt,
	}, nil

}

type CreateAccountUseCaseInput struct {
	id        string
	name      string
	cpf       string
	secret    string
	balance   int
	createdAt *time.Time
}

func NewCreateAccountUseCaseInput(ID, name, cpf, secret string, balance int, createdAt time.Time) *CreateAccountUseCaseInput {
	return &CreateAccountUseCaseInput{
		id:        ID,
		name:      name,
		cpf:       cpf,
		secret:    secret,
		balance:   balance,
		createdAt: &createdAt,
	}
}

type CreateAccountUseCaseOutput struct {
	ID        string
	Name      string
	Balance   int
	CreatedAt *time.Time
}
