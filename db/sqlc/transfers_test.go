package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//for test isolation we will create a new random transaction for each test
func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	debit := account1.Balance - 1

	arg := CreateTransferParams {
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: debit,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)


	return transfer

}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	Transfer1 := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), Transfer1.ID)
	require.NoError(t, err)

	//verify the transfer is deleted
	transfer2, err := testQueries.GetTransfer(context.Background(), Transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	// make sure we have at least 10 records
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	// make sure we return exactly 5 records 
	arg := ListTransfersParams {
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	// verify they are not empty 
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}