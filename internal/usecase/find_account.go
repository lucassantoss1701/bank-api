package usecase

import (
	"lucassantoss1701/bank/internal/entity"
	"time"
)

type FindAccountUseCase struct {
	repostiory entity.AccountRepositoryInterface
}

func NewFindAccountUseCase(repostiory entity.AccountRepositoryInterface) *FindAccountUseCase {
	return &FindAccountUseCase{
		repostiory: repostiory,
	}
}

func (f *FindAccountUseCase) Execute(input *FindAccountUseCaseInput) ([]FindAccountUseCaseOutput, error) {
	accounts, err := f.repostiory.Find(input.limit, input.offset)
	if err != nil {
		return nil, err
	}

	var output []FindAccountUseCaseOutput

	for _, account := range accounts {
		accountOutput := FindAccountUseCaseOutput{
			ID:        account.ID,
			Name:      account.Name,
			CreatedAt: account.CreatedAt,
		}

		output = append(output, accountOutput)
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
	ID        string
	Name      string
	CreatedAt *time.Time
}