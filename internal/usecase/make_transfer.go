package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"
)

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

	originAccount, err := m.accountRepository.FindByID(ctx, input.originAccountID)
	if err != nil {
		return nil, err
	}

	destinationAccount, err := m.accountRepository.FindByID(ctx, input.destinationAccountID)
	if err != nil {
		return nil, err
	}

	transfer, err := entity.NewTransfer(input.id, &originAccount, &destinationAccount, input.amount, input.createdAt)
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

	output := &MakeTransferUseCaseOutput{
		ID:        createdTransfer.ID,
		Amount:    createdTransfer.Amount,
		CreatedAt: createdTransfer.CreatedAt,
		Account: MakeTransferUseCaseAccount{
			ID:        originAccount.ID,
			Name:      originAccount.Name,
			CreatedAt: originAccount.CreatedAt,
		},
	}

	return output, nil
}

type MakeTransferUseCaseInput struct {
	id                   string
	originAccountID      string
	destinationAccountID string
	amount               int
	createdAt            *time.Time
}

func NewMakeTransferUseCaseInput(id string, originAccountID string, destinationAccountID string, amount int, createdAt *time.Time) *MakeTransferUseCaseInput {
	return &MakeTransferUseCaseInput{
		id:                   id,
		originAccountID:      originAccountID,
		destinationAccountID: destinationAccountID,
		amount:               amount,
		createdAt:            createdAt,
	}
}

type MakeTransferUseCaseOutput struct {
	ID        string
	Amount    int
	Account   MakeTransferUseCaseAccount
	CreatedAt *time.Time
}

type MakeTransferUseCaseAccount struct {
	ID        string
	Name      string
	CreatedAt *time.Time
}
