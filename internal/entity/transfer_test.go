package entity_test

import (
	"lucassantoss1701/bank/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func GetBaseOriginAccount(t *testing.T) *entity.Account {
	originAccountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
	originAccountName := "lucas"
	originAccountCPF := "52849254088"
	originAccountsecret := "4578405"
	originAccountBalance := 100
	originAccountCreatedAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
	originAccount, err := entity.NewAccount(originAccountID, originAccountName, originAccountCPF, originAccountsecret, originAccountBalance, &originAccountCreatedAt)

	assert.Nil(t, err)
	assert.NotNil(t, originAccount)
	return originAccount
}

func GetBaseDestinationAccount(t *testing.T) *entity.Account {
	destionationAccountID := "d18551d3-cf13-49ec-b1dc-741a1f8715f6"
	destionationAccountName := "joao"
	destionationAccountCPF := "35768297090"
	destionationAccountSecret := "744637"
	destionationAccountBalance := 200
	destionationAccountCreatedAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)
	destinationAccount, err := entity.NewAccount(destionationAccountID, destionationAccountName, destionationAccountCPF, destionationAccountSecret, destionationAccountBalance, &destionationAccountCreatedAt)

	assert.Nil(t, err)
	assert.NotNil(t, destinationAccount)
	return destinationAccount
}

func TestTransfer_NewTranser(t *testing.T) {
	t.Run("Testing NewTranser when return a valid transfer", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		ownerID := originAccount.ID
		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, destinationAccount, amount, &transferCreatedAt)

		assert.Nil(t, err)
		assert.NotNil(t, transfer)
	})

	t.Run("Testing invalid transfer (empty ID)", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		ownerID := originAccount.ID
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)
		amount := 50

		transfer, err := entity.NewTransfer("", ownerID, originAccount, destinationAccount, amount, &transferCreatedAt)

		assert.Nil(t, transfer)
		assert.NotNil(t, err)
		assert.Equal(t, "ID cannot be empty", err.Error())
	})

	t.Run("Testing invalid transfer (empty ownerID)", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)
		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"
		ownerID := ""
		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, destinationAccount, amount, &transferCreatedAt)
		assert.Nil(t, transfer)
		assert.NotNil(t, err)
		assert.Equal(t, "owner cannot be empty", err.Error())
	})

	t.Run("Testing invalid transfer (nil originAccount)", func(t *testing.T) {

		destinationAccount := GetBaseDestinationAccount(t)
		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"
		ownerID := "135f45fd-8df7-42c5-98cd-84568a8f124c"
		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, nil, destinationAccount, amount, &transferCreatedAt)

		assert.Nil(t, transfer)
		assert.NotNil(t, err)
		assert.Equal(t, "originAccount cannot be nil", err.Error())
	})

	t.Run("Testing invalid transfer (nil destinationAccount)", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"
		ownerID := originAccount.ID
		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, nil, amount, &transferCreatedAt)
		assert.Nil(t, transfer)
		assert.NotNil(t, err)
		assert.Equal(t, "destinationAccount cannot be nil", err.Error())
	})

	t.Run("Testing invalid transfer (negative amount)", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		ownerID := originAccount.ID
		amount := -100
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, destinationAccount, amount, &transferCreatedAt)
		assert.Nil(t, transfer)
		assert.NotNil(t, err)
		assert.Equal(t, "amount cannot be minor than zero", err.Error())
	})

	t.Run("Testing invalid transfer (nil createdAt)", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		ownerID := originAccount.ID
		amount := 100

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, destinationAccount, amount, nil)
		assert.Nil(t, transfer)
		assert.NotNil(t, err)
		assert.Equal(t, "created at cannot be nil", err.Error())
	})
}

func TestTransfer_CanTransfer(t *testing.T) {
	t.Run("Testing CanTransfer when transfer can be performed", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		ownerID := originAccount.ID
		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, destinationAccount, amount, &transferCreatedAt)

		assert.Nil(t, err)
		assert.NotNil(t, transfer)

		err = transfer.CanTransfer()
		assert.Nil(t, err)

	})

	t.Run("Testing CanTransfer when transfer cannot be performed (origin account does not have enough balance)", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		originAccount.Balance = 49

		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		ownerID := originAccount.ID
		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, destinationAccount, amount, &transferCreatedAt)

		assert.Nil(t, err)
		assert.NotNil(t, transfer)

		err = transfer.CanTransfer()
		assert.NotNil(t, err)
		assert.Equal(t, "origin account does not have enough balance", err.Error())
	})

	t.Run("Testing CanTransfer when transfer cannot be performed (origin account does not have enough balance)", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		originAccount.Balance = 49

		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		ownerID := originAccount.ID
		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, destinationAccount, amount, &transferCreatedAt)

		assert.Nil(t, err)
		assert.NotNil(t, transfer)

		err = transfer.CanTransfer()
		assert.NotNil(t, err)
		assert.Equal(t, "origin account does not have enough balance", err.Error())

	})

	t.Run("Testing CanTransfer when transfer cannot be performed (owner ID cannot be different from ID origin account)", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		ownerID := "6f57a36b-dcba-45d0-9456-5e07eeb2dad2"
		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, destinationAccount, amount, &transferCreatedAt)

		assert.Nil(t, err)
		assert.NotNil(t, transfer)

		err = transfer.CanTransfer()
		assert.NotNil(t, err)
		assert.Equal(t, "owner ID cannot be different from ID origin account", err.Error())

	})
}

func TestTransfer_MakeTransfer(t *testing.T) {
	t.Run("Testing MakeTransfer when transfer can be performed with success", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		ownerID := originAccount.ID
		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, destinationAccount, amount, &transferCreatedAt)

		assert.Nil(t, err)
		assert.NotNil(t, transfer)

		expectedBalanceAfterTransferOfOriginAccount := originAccount.Balance - 50
		expectedBalanceAfterTransferOfDestinationAccount := originAccount.Balance + 50

		err = transfer.MakeTransfer()
		assert.Nil(t, err)

		assert.Equal(t, expectedBalanceAfterTransferOfOriginAccount, originAccount.Balance)
		assert.Equal(t, expectedBalanceAfterTransferOfDestinationAccount, destinationAccount.Balance)
	})

	t.Run("Testing MakeTransfer when transfer cannot be performed with success(error on update balance of origin account)", func(t *testing.T) {
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		ownerID := originAccount.ID
		amount := 5000
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, ownerID, originAccount, destinationAccount, amount, &transferCreatedAt)

		assert.Nil(t, err)
		assert.NotNil(t, transfer)

		err = transfer.MakeTransfer()
		assert.NotNil(t, err)
		assert.Equal(t, "error on update balance of origin account: new balance cannot be minor than 0", err.Error())
	})

}
