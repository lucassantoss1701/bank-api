package entity

import (
	"fmt"
	"time"
)

type Transfer struct {
	ID                 string
	OriginAccount      *Account
	DestinationAccount *Account
	Amount             int
	CreatedAt          *time.Time
}

func NewTransfer(ID string, originAccount *Account, destinationAccount *Account, amount int, createdAt *time.Time) (*Transfer, error) {

	if ID == "" {
		ID = NewUUID()
	}

	transfer := &Transfer{
		ID:                 ID,
		OriginAccount:      originAccount,
		DestinationAccount: destinationAccount,
		Amount:             amount,
		CreatedAt:          createdAt,
	}

	err := transfer.isValid()
	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (t *Transfer) isValid() error {
	validationError := NewErrorHandler(ENTITY_ERROR)

	if t.OriginAccount == nil {
		validationError.Add("originAccount cannot be nil")
	}

	if t.DestinationAccount == nil {
		validationError.Add("destinationAccount cannot be nil")
	}

	if t.Amount < 0 {
		validationError.Add("amount cannot be minor than zero")
	}

	if t.CreatedAt == nil {
		validationError.Add("created at cannot be nil")
	}

	if len(validationError.Messages) > 0 {
		return validationError
	}

	return nil
}

func (t *Transfer) MakeTransfer() error {
	validationError := NewErrorHandler(ENTITY_ERROR)

	if t.OriginAccount.ID == t.DestinationAccount.ID {
		validationError.Add("origin account id must be different to destination account id")
		return validationError
	}

	err := t.OriginAccount.removeFromBalance(t.Amount)
	if err != nil {
		validationError.Add(fmt.Sprintf("error on update balance of origin account: %s", err.Error()))
		return validationError
	}

	t.DestinationAccount.addFromBalance(t.Amount)

	return nil
}
