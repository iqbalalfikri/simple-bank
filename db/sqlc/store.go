package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{Queries: New(db), db: db}
}

func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
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

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
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

		err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
			Balance: -arg.Amount,
			ID:      arg.FromAccountID,
		})

		if err != nil {
			return err
		}
		result.FromAccount, err = queries.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
			Balance: arg.Amount,
			ID:      arg.ToAccountID,
		})

		if err != nil {
			return err
		}
		result.ToAccount, err = queries.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}
