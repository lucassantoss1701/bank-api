package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"
)

type ICreateAccountUseCase interface {
	Execute(ctx context.Context, input *CreateAccountUseCaseInput) (*CreateAccountUseCaseOutput, error)
}

type CreateAccountUseCase struct {
	repostiory entity.AccountRepository
}

func NewCreateAccountUseCase(repostiory entity.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		repostiory: repostiory,
	}
}

func (c *CreateAccountUseCase) Execute(ctx context.Context, input *CreateAccountUseCaseInput) (*CreateAccountUseCaseOutput, error) {
	account, err := entity.NewAccount(input.ID, input.Name, input.CPF, input.Secret, input.Balance, input.CreatedAt)
	if err != nil {
		return nil, err
	}

	createdAccount, err := c.repostiory.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return NewCreateAccountUseCaseOutput(createdAccount.ID, createdAccount.Name, createdAccount.Balance, createdAccount.CreatedAt), nil
}

type CreateAccountUseCaseInput struct {
	ID        string     `json:"-"`
	Name      string     `json:"name"`
	CPF       string     `json:"cpf"`
	Secret    string     `json:"secret"`
	Balance   int        `json:"balance"`
	CreatedAt *time.Time `json:"-"`
}

func NewCreateAccountUseCaseInput(ID, name, CPF, secret string, balance int, createdAt time.Time) *CreateAccountUseCaseInput {
	return &CreateAccountUseCaseInput{
		ID:        ID,
		Name:      name,
		CPF:       CPF,
		Secret:    secret,
		Balance:   balance,
		CreatedAt: &createdAt,
	}
}

type CreateAccountUseCaseOutput struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Balance   int        `json:"balance"`
	CreatedAt *time.Time `json:"created_at"`
}

func NewCreateAccountUseCaseOutput(ID string, name string, balance int, createdAt *time.Time) *CreateAccountUseCaseOutput {
	return &CreateAccountUseCaseOutput{
		ID:        ID,
		Name:      name,
		Balance:   balance,
		CreatedAt: createdAt,
	}
}
