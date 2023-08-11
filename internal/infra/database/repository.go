package database

import (
	"context"
	"database/sql"
	"lucassantoss1701/bank/internal/entity"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Db: db,
	}
}

func (r *Repository) BeginTx(ctx context.Context) (entity.TransactionHandler, error) {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *Repository) CommitTx(tx entity.TransactionHandler) error {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		err := sqlTx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) RollbackTx(tx entity.TransactionHandler) error {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		err := sqlTx.Rollback()
		if err != nil {
			return err
		}
	}
	return nil
}
