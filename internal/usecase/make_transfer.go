package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"
)

type IMakeTransferUseCase interface {
	Execute(ctx context.Context, input *MakeTransferUseCaseInput) (*MakeTransferUseCaseOutput, error)
}

type MakeTransferUseCase struct {
	accountRepository  entity.AccountRepository
	transferRepository entity.TransferRepository
	entity.Repository
}

func NewMakeTransferUseCase(accountRepository entity.AccountRepository, transferRepository entity.TransferRepository, repository entity.Repository) *MakeTransferUseCase {
	return &MakeTransferUseCase{
		accountRepository:  accountRepository,
		transferRepository: transferRepository,
		Repository:         repository,
	}
}

func (m *MakeTransferUseCase) Execute(ctx context.Context, input *MakeTransferUseCaseInput) (*MakeTransferUseCaseOutput, error) {

	originAccount, err := m.accountRepository.FindByID(ctx, input.OriginAccount.ID)
	if err != nil {
		return nil, err
	}

	destinationAccount, err := m.accountRepository.FindByID(ctx, input.DestinationAccount.ID)
	if err != nil {
		return nil, err
	}

	transfer, err := entity.NewTransfer(input.ID, &originAccount, &destinationAccount, input.Amount, input.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = transfer.MakeTransfer()
	if err != nil {
		return nil, err
	}

	transaction, err := m.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = m.RollbackTx(transaction)
			panic(r)
		}
		if err != nil {
			_ = m.RollbackTx(transaction)
		} else {
			_ = m.CommitTx(transaction)
		}
	}()

	createdTransfer, err := m.transferRepository.Create(ctx, transfer, transaction)
	if err != nil {
		return nil, err
	}

	originAccount, err = m.accountRepository.UpdateBalance(ctx, transfer.OriginAccount.ID, transfer.OriginAccount.Balance, transaction)
	if err != nil {
		return nil, err
	}

	destinationAccount, err = m.accountRepository.UpdateBalance(ctx, transfer.DestinationAccount.ID, transfer.DestinationAccount.Balance, transaction)
	if err != nil {
		return nil, err
	}

	createdTransfer.OriginAccount = &originAccount
	createdTransfer.DestinationAccount = &destinationAccount

	output := NewMakeTransferUseCaseOutput(&createdTransfer)

	return output, nil
}

type MakeTransferUseCaseInput struct {
	ID                 string                          `json:"-"`
	OriginAccount      MakeTransferUseCaseAccountInput `json:"-"`
	DestinationAccount MakeTransferUseCaseAccountInput `json:"destination_account"`
	Amount             int                             `json:"amount"`
	CreatedAt          *time.Time                      `json:"-"`
}

type MakeTransferUseCaseAccountInput struct {
	ID string
}

func NewMakeTransferUseCaseInput(ID string, originAccountID string, destinationAccountID string, amount int, createdAt *time.Time) *MakeTransferUseCaseInput {
	return &MakeTransferUseCaseInput{
		ID: ID,
		OriginAccount: MakeTransferUseCaseAccountInput{
			ID: originAccountID,
		},
		DestinationAccount: MakeTransferUseCaseAccountInput{
			ID: destinationAccountID,
		},
		Amount:    amount,
		CreatedAt: createdAt,
	}
}

type MakeTransferUseCaseOutput struct {
	ID                 string                     `json:"id"`
	Amount             int                        `json:"amount"`
	OriginAccount      MakeTransferUseCaseAccount `json:"origin_account"`
	DestinationAccount MakeTransferUseCaseAccount `json:"destination_account"`
	CreatedAt          string                     `json:"created_at"`
}

type MakeTransferUseCaseAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewMakeTransferUseCaseOutput(transfer *entity.Transfer) *MakeTransferUseCaseOutput {
	return &MakeTransferUseCaseOutput{
		ID:        transfer.ID,
		Amount:    transfer.Amount,
		CreatedAt: transfer.CreatedAt.Format(time.RFC3339),
		OriginAccount: MakeTransferUseCaseAccount{
			ID:   transfer.OriginAccount.ID,
			Name: transfer.OriginAccount.Name,
		},
		DestinationAccount: MakeTransferUseCaseAccount{
			ID:   transfer.DestinationAccount.ID,
			Name: transfer.DestinationAccount.Name,
		},
	}
}
