package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"
)

type FindTransferencesByAccountUseCase struct {
	repostiory entity.TransferRepository
}

func NewFindTransferencesByAccountUseCase(repostiory entity.TransferRepository) *FindTransferencesByAccountUseCase {
	return &FindTransferencesByAccountUseCase{
		repostiory: repostiory,
	}
}

func (f *FindTransferencesByAccountUseCase) Execute(ctx context.Context, input *FindTransferencesByAccountUseCaseInput) ([]FindTransferencesByAccountUseCaseOutput, error) {
	transfererences, err := f.repostiory.FindByAccountID(ctx, input.accountID, input.limit, input.offset)
	if err != nil {
		return nil, err
	}

	var output []FindTransferencesByAccountUseCaseOutput
	for _, transfer := range transfererences {
		transferOutput := &FindTransferencesByAccountUseCaseOutput{
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

type FindTransferencesByAccountUseCaseInput struct {
	accountID string
	limit     int
	offset    int
}

func NewFindTransferencesByAccountUseCaseInput(accountID string, limit int, offset int) *FindTransferencesByAccountUseCaseInput {
	return &FindTransferencesByAccountUseCaseInput{
		accountID: accountID,
		limit:     limit,
		offset:    offset,
	}
}

type FindTransferencesByAccountUseCaseOutput struct {
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
