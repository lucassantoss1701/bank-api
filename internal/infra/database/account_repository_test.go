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

func GetSQLFindAccounts() string {
	return "SELECT id, name, balance FROM account LIMIT 10 OFFSET 0"
}

func GetSQLFindAccountByID() string {
	return "SELECT id, name, balance FROM account WHERE id = ?"
}

func GetSQLFindByCPF() string {
	return "SELECT id, secret FROM account WHERE cpf = ?"
}

func GetSQLInsertAccount() string {
	return regexp.QuoteMeta("INSERT INTO account (id, name, cpf, secret, balance, created_at) VALUES (?, ?, ?, ?, ?, ?)")
}

func GetSQLUpdateBalanceAccount() string {
	return regexp.QuoteMeta("UPDATE account SET balance = ? WHERE id = ?")
}

func TestAccountRepository_Find(t *testing.T) {
	t.Run("Testing Find when returns two accounts", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		rows := sqlmock.NewRows([]string{"id", "name", "balance"}).
			AddRow("2bd765a6-47bd-4731-9eb2-1e65542f4477", "Lucas", 100).
			AddRow("d18551d3-cf13-49ec-b1dc-741a1f8715f6", "Roger", 200)

		mock.ExpectQuery(GetSQLFindAccounts()).WillReturnRows(rows)

		accounts, err := accountRepository.Find(context.Background(), 10, 0)
		assert.Nil(t, err)
		assert.Len(t, accounts, 2)

		assert.Equal(t, accounts[0].ID, "2bd765a6-47bd-4731-9eb2-1e65542f4477")
		assert.Equal(t, accounts[0].Name, "Lucas")
		assert.Equal(t, accounts[0].Balance, 100)
		assert.Equal(t, accounts[0].CPF, "")
		assert.Equal(t, accounts[0].Secret, "")

		assert.Equal(t, accounts[1].ID, "d18551d3-cf13-49ec-b1dc-741a1f8715f6")
		assert.Equal(t, accounts[1].Name, "Roger")
		assert.Equal(t, accounts[1].Balance, 200)
		assert.Equal(t, accounts[1].CPF, "")
		assert.Equal(t, accounts[1].Secret, "")
	})

	t.Run("Testing Find when returns two accounts and limit is zero", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		rows := sqlmock.NewRows([]string{"id", "name", "balance"}).
			AddRow("2bd765a6-47bd-4731-9eb2-1e65542f4477", "Lucas", 100).
			AddRow("d18551d3-cf13-49ec-b1dc-741a1f8715f6", "Roger", 200)

		mock.ExpectQuery(GetSQLFindAccounts()).WillReturnRows(rows)

		accounts, err := accountRepository.Find(context.Background(), 0, 0)
		assert.Nil(t, err)
		assert.Len(t, accounts, 2)

		assert.Equal(t, accounts[0].ID, "2bd765a6-47bd-4731-9eb2-1e65542f4477")
		assert.Equal(t, accounts[0].Name, "Lucas")
		assert.Equal(t, accounts[0].Balance, 100)
		assert.Equal(t, accounts[0].CPF, "")
		assert.Equal(t, accounts[0].Secret, "")

		assert.Equal(t, accounts[1].ID, "d18551d3-cf13-49ec-b1dc-741a1f8715f6")
		assert.Equal(t, accounts[1].Name, "Roger")
		assert.Equal(t, accounts[1].Balance, 200)
		assert.Equal(t, accounts[1].CPF, "")
		assert.Equal(t, accounts[1].Secret, "")
	})

	t.Run("Testing Find when execute query returns an error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		mock.ExpectQuery(GetSQLFindAccounts()).WillReturnError(errors.New("connection closed"))

		accounts, err := accountRepository.Find(context.Background(), 10, 0)
		assert.NotNil(t, err)
		assert.Equal(t, "connection closed", err.Error())
		assert.Len(t, accounts, 0)
	})

}

func TestAccountRepository_FindByID(t *testing.T) {

	t.Run("Testing findByID when returns one account", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		rows := sqlmock.NewRows([]string{"id", "name", "balance"}).
			AddRow("2bd765a6-47bd-4731-9eb2-1e65542f4477", "Lucas", 100)

		mock.ExpectQuery(GetSQLFindAccountByID()).WithArgs("2bd765a6-47bd-4731-9eb2-1e65542f4477").WillReturnRows(rows)

		account, err := accountRepository.FindByID(context.Background(), "2bd765a6-47bd-4731-9eb2-1e65542f4477")
		assert.Nil(t, err)
		assert.Equal(t, "2bd765a6-47bd-4731-9eb2-1e65542f4477", account.ID)
		assert.Equal(t, "Lucas", account.Name)
		assert.Equal(t, 100, account.Balance)
	})

	t.Run("Testing FindByID when execute returns an error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		mock.ExpectQuery(GetSQLFindAccountByID()).WithArgs("2bd765a6-47bd-4731-9eb2-1e65542f4477").WillReturnError(errors.New("connection closed"))

		account, err := accountRepository.FindByID(context.Background(), "2bd765a6-47bd-4731-9eb2-1e65542f4477")
		assert.NotNil(t, err)
		assert.Equal(t, "connection closed", err.Error())
		assert.Empty(t, account.ID)

	})

	t.Run("Testing findByID when scan returns an error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		rows := sqlmock.NewRows([]string{"id", "name", "balance"}).
			AddRow("2bd765a6-47bd-4731-9eb2-1e65542f4477", "Lucas", 100).CloseError(errors.New("error on scan"))

		mock.ExpectQuery(GetSQLFindAccountByID()).WithArgs("2bd765a6-47bd-4731-9eb2-1e65542f4477").WillReturnRows(rows)

		account, err := accountRepository.FindByID(context.Background(), "2bd765a6-47bd-4731-9eb2-1e65542f4477")
		assert.NotNil(t, err)
		assert.Equal(t, "error on scan", err.Error())
		assert.Empty(t, account.ID)
	})

	t.Run("Testing findByID when scan returns an error (no rows)", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		rows := sqlmock.NewRows([]string{"id", "name", "balance"}).
			AddRow("2bd765a6-47bd-4731-9eb2-1e65542f4477", "Lucas", 100).CloseError(errors.New("sql: no rows in result set"))

		mock.ExpectQuery(GetSQLFindAccountByID()).WithArgs("2bd765a6-47bd-4731-9eb2-1e65542f4477").WillReturnRows(rows)

		account, err := accountRepository.FindByID(context.Background(), "2bd765a6-47bd-4731-9eb2-1e65542f4477")
		assert.NotNil(t, err)
		assert.Equal(t, "sql: no rows in result set", err.Error())
		assert.Empty(t, account.ID)
	})

}

func TestAccountRepository_Create(t *testing.T) {
	t.Run("Testing Create when account is create with successful", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		createdAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)

		account := &entity.Account{
			ID:        "2bd765a6-47bd-4731-9eb2-1e65542f4477",
			Name:      "John",
			CPF:       "00634020099",
			Secret:    "4578405",
			Balance:   200,
			CreatedAt: &createdAt,
		}

		mock.ExpectExec(GetSQLInsertAccount()).
			WithArgs(account.ID, account.Name, account.CPF, account.Secret, account.Balance, account.CreatedAt).
			WillReturnResult(sqlmock.NewResult(0, 1))

		createdAccount, err := accountRepository.Create(context.Background(), account)
		assert.Nil(t, err)
		assert.Equal(t, account.ID, createdAccount.ID)
		assert.Equal(t, account.Name, createdAccount.Name)
		assert.Equal(t, account.Balance, createdAccount.Balance)
	})

	t.Run("Testing Create when ExecContext returns an erros", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		createdAt := time.Date(2023, 8, 5, 8, 22, 00, 00, time.UTC)

		account := &entity.Account{
			ID:        "2bd765a6-47bd-4731-9eb2-1e65542f4477",
			Name:      "John",
			CPF:       "00634020099",
			Secret:    "4578405",
			Balance:   200,
			CreatedAt: &createdAt,
		}

		mock.ExpectExec(GetSQLInsertAccount()).
			WithArgs(account.ID, account.Name, account.CPF, account.Secret, account.Balance, account.CreatedAt).
			WillReturnError(errors.New("connection closed"))

		createdAccount, err := accountRepository.Create(context.Background(), account)
		assert.NotNil(t, err)
		assert.Equal(t, "connection closed", err.Error())
		assert.Empty(t, createdAccount.ID)
	})
}

func TestAccountRepository_UpdateBalance(t *testing.T) {
	t.Run("Testing UpdateBalance when successful", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		accountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
		newBalance := 150
		accountName := "Lucas"

		mock.ExpectExec(GetSQLUpdateBalanceAccount()).
			WithArgs(newBalance, accountID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		rows := sqlmock.NewRows([]string{"id", "name", "balance"}).
			AddRow(accountID, accountName, newBalance)

		mock.ExpectQuery(GetSQLFindAccountByID()).WithArgs(accountID).WillReturnRows(rows)

		updatedAccount, err := accountRepository.UpdateBalance(context.Background(), accountID, newBalance)
		assert.Nil(t, err)
		assert.Equal(t, accountID, updatedAccount.ID)
		assert.Equal(t, accountName, updatedAccount.Name)
		assert.Equal(t, newBalance, updatedAccount.Balance)
	})

	t.Run("Testing UpdateBalance when ExecContext returns an error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		accountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
		newBalance := 150
		accountName := "Lucas"

		mock.ExpectExec(GetSQLUpdateBalanceAccount()).
			WithArgs(newBalance, accountID).
			WillReturnError(errors.New("error on update balance"))

		rows := sqlmock.NewRows([]string{"id", "name", "balance"}).
			AddRow(accountID, accountName, newBalance)

		mock.ExpectQuery(GetSQLFindAccountByID()).WithArgs(accountID).WillReturnRows(rows)

		updatedAccount, err := accountRepository.UpdateBalance(context.Background(), accountID, newBalance)
		assert.NotNil(t, err)
		assert.Equal(t, "error on update balance", err.Error())
		assert.Empty(t, updatedAccount.ID)
	})

	t.Run("Testing UpdateBalance when successful with transaction", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		accountID := "2bd765a6-47bd-4731-9eb2-1e65542f4477"
		newBalance := 150
		accountName := "Lucas"

		mock.ExpectExec(GetSQLUpdateBalanceAccount()).
			WithArgs(newBalance, accountID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		rows := sqlmock.NewRows([]string{"id", "name", "balance"}).
			AddRow(accountID, accountName, newBalance)

		mock.ExpectQuery(GetSQLFindAccountByID()).WithArgs(accountID).WillReturnRows(rows)

		updatedAccount, err := accountRepository.UpdateBalance(context.Background(), accountID, newBalance, db)
		assert.Nil(t, err)
		assert.Equal(t, accountID, updatedAccount.ID)
		assert.Equal(t, accountName, updatedAccount.Name)
		assert.Equal(t, newBalance, updatedAccount.Balance)
	})
}

func TestAccountRepository_FindByCPF(t *testing.T) {

	t.Run("Testing FindByCPF when returns one account", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		rows := sqlmock.NewRows([]string{"id", "secret"}).
			AddRow("2bd765a6-47bd-4731-9eb2-1e65542f4477", "secret")

		mock.ExpectQuery(GetSQLFindByCPF()).WithArgs("35768297090").WillReturnRows(rows)

		account, err := accountRepository.FindByCPF(context.Background(), "35768297090")
		assert.Nil(t, err)
		assert.Equal(t, "2bd765a6-47bd-4731-9eb2-1e65542f4477", account.ID)
		assert.Equal(t, "secret", account.Secret)
	})

	t.Run("Testing FindByID when execute returns an error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		mock.ExpectQuery(GetSQLFindByCPF()).WithArgs("35768297090").WillReturnError(errors.New("connection closed"))

		account, err := accountRepository.FindByCPF(context.Background(), "35768297090")
		assert.NotNil(t, err)
		assert.Equal(t, "connection closed", err.Error())
		assert.Empty(t, account.ID)

	})

	t.Run("Testing findByID when scan returns an error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		rows := sqlmock.NewRows([]string{"id", "secret"}).
			AddRow("2bd765a6-47bd-4731-9eb2-1e65542f4477", "secret").CloseError(errors.New("error on scan"))

		mock.ExpectQuery(GetSQLFindByCPF()).WithArgs("35768297090").WillReturnRows(rows)

		account, err := accountRepository.FindByCPF(context.Background(), "35768297090")
		assert.NotNil(t, err)
		assert.Equal(t, "error on scan", err.Error())
		assert.Empty(t, account.ID)
	})

	t.Run("Testing findByID when scan returns an error (no rows)", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer db.Close()

		accountRepository := database.NewAccountRepository(db)

		rows := sqlmock.NewRows([]string{"id", "secret"}).
			AddRow("2bd765a6-47bd-4731-9eb2-1e65542f4477", "secret").CloseError(errors.New("sql: no rows in result set"))

		mock.ExpectQuery(GetSQLFindByCPF()).WithArgs("35768297090").WillReturnRows(rows)

		account, err := accountRepository.FindByCPF(context.Background(), "35768297090")
		assert.NotNil(t, err)
		assert.Equal(t, "sql: no rows in result set", err.Error())
		assert.Empty(t, account.ID)
	})

}
