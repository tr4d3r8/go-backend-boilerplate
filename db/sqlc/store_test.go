package db

import (
	"context"
	"testing"
	"github.com/stretchr/testify/require"
)


func TestTranferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions 
	n := 5
	amount := int64(10)

	errs := make(chan error) //chan errors used to connect concurrent go routines 
	results := make(chan TransferTxResult) // create channel 

	for i := 0; i <n; i++ {
		go func ()  { // cannot use testify as this is running its own routine 
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			// channel on left data to send on right 
			errs <- err
			results <- result

		}()
	}
	// check results from outside the routine 
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		
		result := <-results
		require.NotEmpty(t, result)
		
		// check transfer 
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries of the result FROM
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//check entries of the result TO
		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, account2.ID, ToEntry.AccountID)
		require.Equal(t, amount, ToEntry.Amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), ToEntry.ID)
		require.NoError(t, err)

		// TODO: check accounts balances once implimented... 
	}
}