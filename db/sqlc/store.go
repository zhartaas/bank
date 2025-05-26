package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := store.Queries.WithTx(tx)

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

type TransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferParams) (TransferResult, error) {
	result := &TransferResult{}

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: update accounts balance
		// goroutine safe

		// approach to handle postgres lock. Locks in postgresql and can be solved with FOR UPDATE or any else
		// good showcase to handle postgres lock.

		//account1, err := q.GetAccount(context.Background(), arg.FromAccountID)
		//if err != nil {
		//	return err
		//}
		//result.FromAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
		//	ID:      arg.FromAccountID,
		//	Balance: account1.Balance - arg.Amount,
		//})
		//if err != nil {
		//	return nil
		//}
		//account2, err := q.GetAccount(context.Background(), arg.ToAccountID)
		//if err != nil {
		//	return err
		//}
		//result.ToAccount, err := q.UpdateAccount(context.Background(), UpdateAccountParams{
		//	ID: arg.ToAccountID,
		//	Balance: account2.Balance + arg.Amount,
		//})
		result.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return nil
		}
		result.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})

		return nil
	})

	return *result, err
}
