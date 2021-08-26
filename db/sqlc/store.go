package db

import (
	"context"
	"database/sql"
	"fmt"
)

// provides all functions to execute db queries and transactions 
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// provide all functions to execute db queries and combinations
// adds queries by composition
type SQLStore struct {
	*Queries
	db *sql.DB
}

// create new store object 
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:		db, 
		Queries: New(db),
	}
}

// executes a function within a database transaction 
// executes or roles back based on return 
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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

var txKey = struct{}{}


// performs money transfer between accounts
// creates transfer record 
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
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
		
		// make sure we always execute transactions in a standard order 
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

	}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context, 
	q *Queries,
	accountID1 int64,
	amount1 int64, 
	accountID2 int64, 
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})
	return
}