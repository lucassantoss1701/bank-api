package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"
)

type FindTransfersByAccountUseCase struct {
	repostiory entity.TransferRepository
}

func NewFindTransfersByAccountUseCase(repostiory entity.TransferRepository) *FindTransfersByAccountUseCase {
	return &FindTransfersByAccountUseCase{
		repostiory: repostiory,
	}
}

func (f *FindTransfersByAccountUseCase) Execute(ctx context.Context, input *FindTransfersByAccountUseCaseInput) ([]FindTransfersByAccountUseCaseOutput, error) {
	transfererences, err := f.repostiory.FindByAccountID(ctx, input.accountID, input.limit, input.offset)
	if err != nil {
		return nil, err
	}

	var output []FindTransfersByAccountUseCaseOutput
	for _, transfer := range transfererences {
		transferOutput := &FindTransfersByAccountUseCaseOutput{
			ID:        transfer.ID,
			Amount:    transfer.Amount,
			CreatedAt: transfer.CreatedAt,
			DestinationAccount: account{
				ID:        transfer.DestinationAccount.ID,
				Name:      transfer.DestinationAccount.Name,
				CreatedAt: transfer.DestinationAccount.CreatedAt,
			},
		}

		output = append(output, *transferOutput)

	}

	return output, nil

}

type FindTransfersByAccountUseCaseInput struct {
	accountID string
	limit     int
	offset    int
}

func NewFindTransfersByAccountUseCaseInput(accountID string, limit int, offset int) *FindTransfersByAccountUseCaseInput {
	return &FindTransfersByAccountUseCaseInput{
		accountID: accountID,
		limit:     limit,
		offset:    offset,
	}
}

type FindTransfersByAccountUseCaseOutput struct {
	ID                 string
	DestinationAccount account
	Amount             int
	CreatedAt          *time.Time
}

type account struct {
	ID        string
	Name      string
	CreatedAt *time.Time
}
