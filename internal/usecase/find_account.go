package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"
)

type IFindAccountUseCase interface {
	Execute(ctx context.Context, input *FindAccountUseCaseInput) ([]FindAccountUseCaseOutput, error)
}

type FindAccountUseCase struct {
	repostiory entity.AccountRepository
}

func NewFindAccountUseCase(repostiory entity.AccountRepository) *FindAccountUseCase {
	return &FindAccountUseCase{
		repostiory: repostiory,
	}
}

func (f *FindAccountUseCase) Execute(ctx context.Context, input *FindAccountUseCaseInput) ([]FindAccountUseCaseOutput, error) {
	accounts, err := f.repostiory.Find(ctx, input.limit, input.offset)
	if err != nil {
		return nil, err
	}

	var output []FindAccountUseCaseOutput

	for _, account := range accounts {

		accountOutput := NewFindAccountUseCaseOutput(account.ID, account.Name, account.CreatedAt)

		output = append(output, *accountOutput)
	}

	return output, nil
}

type FindAccountUseCaseInput struct {
	limit  int
	offset int
}

func NewFindAccountUseCaseInput(limit int, offset int) *FindAccountUseCaseInput {
	return &FindAccountUseCaseInput{
		limit:  limit,
		offset: offset,
	}
}

type FindAccountUseCaseOutput struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
}

func NewFindAccountUseCaseOutput(id string, name string, createdAt *time.Time) *FindAccountUseCaseOutput {
	return &FindAccountUseCaseOutput{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
	}
}
