package entity

import (
	"fmt"
	"time"
)

type Transfer struct {
	ID                 string
	OwnerID            string
	OriginAccount      *Account
	DestinationAccount *Account
	Amount             int
	CreatedAt          *time.Time
}

func NewTransfer(ID string, ownerID string, originAccount *Account, destinationAccount *Account, amount int, createdAt *time.Time) (*Transfer, error) {
	transfer := &Transfer{
		ID:                 ID,
		OwnerID:            ownerID,
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
	validationError := &ValidationError{}

	if t.ID == "" {
		validationError.Add("ID cannot be empty")
	}

	if t.OwnerID == "" {
		validationError.Add("owner cannot be empty")
	}

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

func (t *Transfer) CanTransfer() error {
	validationError := &ValidationError{}
	if !t.OriginAccount.hasSufficientBalance(t.Amount) {
		validationError.Add("origin account does not have enough balance")
	}

	if t.OriginAccount.ID != t.OwnerID {
		validationError.Add("owner ID cannot be different from ID origin account")
	}

	if len(validationError.Messages) > 0 {
		return validationError
	}

	return nil
}

func (t *Transfer) MakeTransfer() error {
	validationError := &ValidationError{}

	err := t.OriginAccount.setBalance(-t.Amount)
	if err != nil {
		validationError.Add(fmt.Sprintf("error on update balance of origin account: %s", err.Error()))
		return validationError
	}

	t.DestinationAccount.setBalance(t.Amount)

	return nil
}
