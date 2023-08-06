package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"
)

type FindTransferByAccountUseCase struct {
	repostiory entity.TransferRepositoryInterface
}

func NewFindTransferByAccountUseCase(repostiory entity.TransferRepositoryInterface) *FindTransferByAccountUseCase {
	return &FindTransferByAccountUseCase{
		repostiory: repostiory,
	}
}

func (f *FindTransferByAccountUseCase) Execute(ctx context.Context, input *FindTransferByAccountUseCaseInput) ([]FindTransferByAccountUseCaseOutput, error) {
	transfererences, err := f.repostiory.FindByAccountID(ctx, input.accountID, input.limit, input.offset)
	if err != nil {
		return nil, err
	}

	var output []FindTransferByAccountUseCaseOutput
	for _, transfer := range transfererences {
		transferOutput := &FindTransferByAccountUseCaseOutput{
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

type FindTransferByAccountUseCaseInput struct {
	accountID string
	limit     int
	offset    int
}

func NewFindTransferByAccountUseCaseInput(accountID string, limit int, offset int) *FindTransferByAccountUseCaseInput {
	return &FindTransferByAccountUseCaseInput{
		accountID: accountID,
		limit:     limit,
		offset:    offset,
	}
}

type FindTransferByAccountUseCaseOutput struct {
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
