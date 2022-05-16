package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{Queries: New(db), db: db}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int32 `json:"from_account_id"`
	ToAccountID   int32 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error

		res, err := queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		id, _ := res.LastInsertId()

		transfer := Transfer{
			ID:            int32(id),
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		}
		result.Transfer = transfer

		res, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		id, _ = res.LastInsertId()
		entry := Entry{
			ID:        int32(id),
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}
		result.FromEntry = entry

		res, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		id, _ = res.LastInsertId()
		entry = Entry{
			ID:        int32(id),
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}
		result.ToEntry = entry

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, queries, int64(arg.FromAccountID), -arg.Amount, int64(arg.ToAccountID), arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, queries, int64(arg.ToAccountID), arg.Amount, int64(arg.FromAccountID), -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	queries *Queries,
	accountID1,
	amount1,
	accountID2,
	amount2 int64,
) (account1, account2 Account, err error) {
	err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
		Balance: amount1,
		ID:      int32(accountID1),
	})
	if err != nil {
		return
	}

	err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
		Balance: amount2,
		ID:      int32(accountID2),
	})
	if err != nil {
		return
	}

	account1, err = queries.GetAccount(ctx, int32(accountID1))
	if err != nil {
		return
	}

	account2, err = queries.GetAccount(ctx, int32(accountID2))
	if err != nil {
		return
	}
	return
}
