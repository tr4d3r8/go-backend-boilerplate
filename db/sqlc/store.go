package db

import (
	"context"
	"database/sql"
	"fmt"
)

// provide all functions to execute db queries and combinations
// adds queries by composition
type Store struct {
	*Queries
	db *sql.DB
}

// create new store object 
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:		db, 
		Queries: New(db),
	}
}

// executes a function within a database transaction 
// executes or roles back based on return 
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// get back new query object 
	q := New(tx)
	err = fn(q)
	// rollback on error - if rollback error return both errors 
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID	int64 	`json:"from_account_id"`
	ToAccountID		int64	`json:"to_account_id`
	Amount			int64	`json:"amount"`
}

// result after transaction 
type TransferTxResult struct {
	Transfer 		Transfer	`json:"transfer"`
	FromAccount		Account		`json:"from_account"`
	ToAccount		Account		`json:"to_account"`
	FromEntry		Entry		`json:"from_entry"`
	ToEntry			Entry		`json:"to_entry"`  
}

// performs money transfer between accounts
// creates transfer record 
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	//create empty result 
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error 
		
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		
		// moving money out of account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		// moving money into an account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		
		//TODO: update account balance 
		// get account -> update its balance -> send it back to db 
		// need to add locking mechanism 

		// this moves money out of the account 
		account1, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: 	arg.FromAccountID,
			Balance: 	account1.Balance - arg.Amount,
		})
		if err != nil {
			return err 
		}

		// this moves money into the account
		account2, err := q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: 	arg.ToAccountID,
			Balance: 	account2.Balance + arg.Amount,
		})
		if err != nil {
			return err 
		}

		return nil
	})

	return result, err
}