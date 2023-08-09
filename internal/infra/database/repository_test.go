package database_test

import (
	"context"
	"database/sql"
	"lucassantoss1701/bank/internal/infra/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_BeginTx(t *testing.T) {
	t.Run("Testing BeginTx when successful", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repository := database.NewRepository(db)

		ctx := context.Background()

		mock.ExpectBegin()
		tx, err := repository.BeginTx(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	})

	t.Run("Testing BeginTx when error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repository := database.NewRepository(db)

		ctx := context.Background()

		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
		tx, err := repository.BeginTx(ctx)
		assert.NotNil(t, err)
		assert.Nil(t, tx)
		assert.Equal(t, sql.ErrConnDone, err)
	})
}

func TestRepository_CommitTx(t *testing.T) {
	t.Run("Testing CommitTx when successful", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repository := database.NewRepository(db)

		mock.ExpectBegin()
		tx, _ := db.Begin()

		mock.ExpectCommit()

		err := repository.CommitTx(tx)
		assert.Nil(t, err)
	})

	t.Run("Testing CommitTx when error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repository := database.NewRepository(db)

		mock.ExpectBegin()
		tx, _ := db.Begin()

		mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

		err := repository.CommitTx(tx)
		assert.NotNil(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})
}

func TestRepository_RollbackTx(t *testing.T) {
	t.Run("Testing RollbackTx when successful", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repository := database.NewRepository(db)

		mock.ExpectBegin()
		tx, _ := db.Begin()

		mock.ExpectRollback()

		err := repository.RollbackTx(tx)
		assert.Nil(t, err)
	})

	t.Run("Testing RollbackTx when error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		repository := database.NewRepository(db)

		mock.ExpectBegin()
		tx, _ := db.Begin()

		mock.ExpectRollback().WillReturnError(sql.ErrConnDone)

		err := repository.RollbackTx(tx)
		assert.NotNil(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})
}
