package database

import (
	"context"
	"database/sql"
	"lucassantoss1701/bank/internal/entity"
)

type TransferRepository struct {
	Db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{
		Db: db,
	}
}

func (r *TransferRepository) FindByAccountID(ctx context.Context, AccountID string, limit, offset int) ([]entity.Transfer, error) {
	query := `
		SELECT t.id, t.amount, t.created_at,
			d.id AS destination_account_id, d.name AS destination_account_name
		FROM transfer t
		INNER JOIN account o ON t.origin_account_id = o.id
		INNER JOIN account d ON t.destination_account_id = d.id
		WHERE t.origin_account_id = ?
		LIMIT ? OFFSET ?
	`

	rows, err := r.Db.QueryContext(ctx, query, AccountID, limit, offset)
	if err != nil {
		return nil, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add(err.Error())
	}
	defer rows.Close()

	transfers := []entity.Transfer{}
	for rows.Next() {
		var transfer entity.Transfer
		var originAccount entity.Account
		var destinationAccount entity.Account

		err := rows.Scan(
			&transfer.ID, &transfer.Amount, &transfer.CreatedAt,
			&destinationAccount.ID, &destinationAccount.Name,
		)
		if err != nil {
			return nil, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add(err.Error())
		}

		transfer.OriginAccount = &originAccount
		transfer.DestinationAccount = &destinationAccount

		transfers = append(transfers, transfer)
	}

	if err = rows.Err(); err != nil {
		return nil, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add(err.Error())
	}

	return transfers, nil
}

func (r *TransferRepository) Create(ctx context.Context, transfer *entity.Transfer, tx ...entity.TransactionHandler) (entity.Transfer, error) {
	var executor entity.TransactionHandler
	if len(tx) > 0 {
		executor = tx[0]
	} else {
		executor = r.Db
	}

	query := `
		INSERT INTO transfer (id, origin_account_id, destination_account_id, amount, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := executor.ExecContext(
		ctx, query, transfer.ID, transfer.OriginAccount.ID, transfer.DestinationAccount.ID, transfer.Amount, transfer.CreatedAt,
	)
	if err != nil {
		return entity.Transfer{}, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add(err.Error())
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return entity.Transfer{}, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add(err.Error())
	}

	if affectedRows != 1 {
		return entity.Transfer{}, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add("unexpected number of affected rows")
	}

	return *transfer, nil
}
