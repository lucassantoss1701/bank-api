package database_test

import (
	"context"
	"errors"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/infra/database"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func GetSQLFindTransfersByAccountID() string {
	return regexp.QuoteMeta(`SELECT t.id, t.amount, t.created_at, d.id AS destination_account_id, d.name AS destination_account_name FROM transfer t INNER JOIN account o ON t.origin_account_id = o.id INNER JOIN account d ON t.destination_account_id = d.id WHERE t.origin_account_id = ? LIMIT ? OFFSET ?`)
}

func GetSQLTransferInsertQuery() string {
	return regexp.QuoteMeta(`INSERT INTO transfer (id, origin_account_id, destination_account_id, amount, created_at) VALUES (?, ?, ?, ?, ?)`)
}

func TestTransferRepository_FindByAccountID(t *testing.T) {
	t.Run("Testing FindByAccountID when successful", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		transferRepository := database.NewTransferRepository(db)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"
		originAccountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"

		destinationAccountID := "d18551d3-cf13-49ec-b1dc-741a1f8715f6"
		destinationAccountName := "Roger"

		limit := 10
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "amount", "created_at",
			"destination_account_id", "destination_account_name",
		}).AddRow(
			transferID, 100, time.Now(),
			destinationAccountID, destinationAccountName,
		)

		mock.ExpectQuery(GetSQLFindTransfersByAccountID()).
			WithArgs(originAccountID, limit, offset).
			WillReturnRows(rows)

		transfers, err := transferRepository.FindByAccountID(context.Background(), originAccountID, limit, offset)
		assert.Nil(t, err)
		assert.Len(t, transfers, 1)
		assert.Equal(t, transferID, transfers[0].ID)
		assert.Equal(t, 100, transfers[0].Amount)

		assert.Equal(t, destinationAccountID, transfers[0].DestinationAccount.ID)
		assert.Equal(t, destinationAccountName, transfers[0].DestinationAccount.Name)
	})

	t.Run("Testing FindByAccountID when QueryContext returns a error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		transferRepository := database.NewTransferRepository(db)

		originAccountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"

		limit := 10
		offset := 0

		mock.ExpectQuery(GetSQLFindTransfersByAccountID()).
			WithArgs(originAccountID, limit, offset).
			WillReturnError(errors.New("connection closed"))

		transfers, err := transferRepository.FindByAccountID(context.Background(), originAccountID, limit, offset)
		assert.NotNil(t, err)
		assert.Equal(t, "connection closed", err.Error())
		assert.Len(t, transfers, 0)
	})

	t.Run("Testing FindByAccountID when rows.Error returns a error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		transferRepository := database.NewTransferRepository(db)

		originAccountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"

		limit := 10
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "amount", "created_at",
			"origin_account_id", "origin_account_name", "origin_account_balance",
			"destination_account_id", "destination_account_name", "destination_account_balance",
		}).CloseError(errors.New("error on scan"))

		mock.ExpectQuery(GetSQLFindTransfersByAccountID()).
			WithArgs(originAccountID, limit, offset).
			WillReturnRows(rows)

		transfers, err := transferRepository.FindByAccountID(context.Background(), originAccountID, limit, offset)
		assert.NotNil(t, err)
		assert.Equal(t, "error on scan", err.Error())
		assert.Len(t, transfers, 0)

	})

}

func TestTransferRepository_Create(t *testing.T) {
	t.Run("Testing Create when successful", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		transferRepository := database.NewTransferRepository(db)

		ctx := context.Background()
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, originAccount, destinationAccount, amount, &transferCreatedAt)
		assert.Nil(t, err)

		mock.ExpectExec(GetSQLTransferInsertQuery()).
			WithArgs(transfer.ID, transfer.OriginAccount.ID, transfer.DestinationAccount.ID, transfer.Amount, transfer.CreatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		createdTransfer, err := transferRepository.Create(ctx, transfer)
		assert.Nil(t, err)
		assert.Equal(t, *transfer, createdTransfer)
	})

	t.Run("Testing Create when successful and use transactionHandler", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		transferRepository := database.NewTransferRepository(db)

		ctx := context.Background()
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, originAccount, destinationAccount, amount, &transferCreatedAt)
		assert.Nil(t, err)

		mock.ExpectExec(GetSQLTransferInsertQuery()).
			WithArgs(transfer.ID, transfer.OriginAccount.ID, transfer.DestinationAccount.ID, transfer.Amount, transfer.CreatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		createdTransfer, err := transferRepository.Create(ctx, transfer, db)
		assert.Nil(t, err)
		assert.Equal(t, *transfer, createdTransfer)
	})

	t.Run("Testing Create when executor.ExecContext returns an error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		transferRepository := database.NewTransferRepository(db)

		ctx := context.Background()
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, originAccount, destinationAccount, amount, &transferCreatedAt)
		assert.Nil(t, err)

		mock.ExpectExec(GetSQLTransferInsertQuery()).
			WithArgs(transfer.ID, transfer.OriginAccount.ID, transfer.DestinationAccount.ID, transfer.Amount, transfer.CreatedAt).
			WillReturnError(errors.New("database error"))

		createdTransfer, err := transferRepository.Create(ctx, transfer)
		assert.NotNil(t, err)
		assert.Equal(t, entity.Transfer{}, createdTransfer)
		assert.Equal(t, "database error", err.Error())
	})

	t.Run("Testing Create when RowsAffected returns unexpected value", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		transferRepository := database.NewTransferRepository(db)

		ctx := context.Background()
		originAccount := GetBaseDestinationAccount(t)
		destinationAccount := GetBaseDestinationAccount(t)

		transferID := "fc84682a-3045-4bdf-b91c-10be19f89452"

		amount := 50
		transferCreatedAt := time.Date(2023, 8, 5, 9, 55, 00, 00, time.UTC)

		transfer, err := entity.NewTransfer(transferID, originAccount, destinationAccount, amount, &transferCreatedAt)
		assert.Nil(t, err)
		mock.ExpectExec(GetSQLTransferInsertQuery()).
			WithArgs(transfer.ID, transfer.OriginAccount.ID, transfer.DestinationAccount.ID, transfer.Amount, transfer.CreatedAt).
			WillReturnResult(sqlmock.NewResult(1, 0))

		createdTransfer, err := transferRepository.Create(ctx, transfer)
		assert.NotNil(t, err)
		assert.Equal(t, entity.Transfer{}, createdTransfer)
		assert.Equal(t, "unexpected number of affected rows", err.Error())
	})
}

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
