package database

import (
	"context"
	"database/sql"
	"fmt"
	"lucassantoss1701/bank/internal/entity"
	"strings"
)

type AccountRepository struct {
	Db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{Db: db}
}

func (r *AccountRepository) Find(ctx context.Context, limit, offset int) ([]entity.Account, error) {
	if limit == 0 {
		limit = 10
	}

	query := fmt.Sprintf("SELECT id, name, balance, created_at FROM account LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add(err.Error())
	}
	defer rows.Close()

	var accounts []entity.Account
	for rows.Next() {
		var account entity.Account
		rows.Scan(&account.ID, &account.Name, &account.Balance, &account.CreatedAt)
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *AccountRepository) FindByID(ctx context.Context, ID string) (entity.Account, error) {
	query := "SELECT id, name, balance FROM account WHERE id = ?"

	row := r.Db.QueryRowContext(ctx, query, ID)

	var account entity.Account
	err := row.Scan(&account.ID, &account.Name, &account.Balance)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return entity.Account{}, entity.NewErrorHandler(entity.NOT_FOUND_ERROR).Add(fmt.Sprintf("not found account: %s", ID))
		}
		return entity.Account{}, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add(err.Error())
	}

	return account, nil
}

func (r *AccountRepository) Create(ctx context.Context, account *entity.Account) (entity.Account, error) {

	query := "INSERT INTO account (id, name, cpf, secret, balance, created_at) VALUES (?, ?, ?, ?, ?, ?)"

	_, err := r.Db.ExecContext(ctx, query, account.ID, account.Name, account.CPF, account.Secret, account.Balance, account.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			return entity.Account{}, entity.NewErrorHandler(entity.CONFLICT_ERROR).Add(err.Error())
		}
		return entity.Account{}, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add(err.Error())
	}

	return *account, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, accountID string, newBalance int, tx ...entity.TransactionHandler) (entity.Account, error) {
	var executor entity.TransactionHandler
	if len(tx) > 0 {
		executor = tx[0]
	} else {
		executor = r.Db
	}

	query := "UPDATE account SET balance = ? WHERE id = ?"

	_, err := executor.ExecContext(ctx, query, newBalance, accountID)
	if err != nil {
		return entity.Account{}, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add(err.Error())
	}

	return r.FindByID(ctx, accountID)
}

func (r *AccountRepository) FindByCPF(ctx context.Context, CPF string) (entity.Account, error) {
	query := "SELECT id, secret FROM account WHERE cpf = ?"

	var account entity.Account

	row := r.Db.QueryRowContext(ctx, query, CPF)
	err := row.Scan(&account.ID, &account.Secret)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return entity.Account{}, entity.NewErrorHandler(entity.NOT_FOUND_ERROR).Add(fmt.Sprintf("not found account by CPF: %s", CPF))
		}
		return entity.Account{}, entity.NewErrorHandler(entity.INTERNAL_ERROR).Add(err.Error())
	}

	return account, nil
}
