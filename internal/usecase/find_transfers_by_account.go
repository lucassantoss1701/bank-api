package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"
)

type IFindTransfersByAccountUseCase interface {
	Execute(ctx context.Context, input *FindTransfersByAccountUseCaseInput) ([]FindTransfersByAccountUseCaseOutput, error)
}

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

	output := []FindTransfersByAccountUseCaseOutput{}
	for _, transfer := range transfererences {
		output = append(output, *NewFindTransfersByAccountUseCaseOutput(transfer))

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
	ID                 string  `json:"id"`
	DestinationAccount account `json:"destination_account"`
	Amount             int     `json:"amount"`
	CreatedAt          string  `json:"created_at"`
}

type account struct {
	ID   string
	Name string
}

func NewFindTransfersByAccountUseCaseOutput(transfer entity.Transfer) *FindTransfersByAccountUseCaseOutput {
	return &FindTransfersByAccountUseCaseOutput{
		ID:        transfer.ID,
		Amount:    transfer.Amount,
		CreatedAt: transfer.CreatedAt.Format(time.RFC3339),
		DestinationAccount: account{
			ID:   transfer.DestinationAccount.ID,
			Name: transfer.DestinationAccount.Name,
		},
	}
}
