package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// for test isolation we will create a random entry for each test
func createRandomEntry(t *testing.T) Entry {
	entry1 := createRandomTransfer(t)

	argTo := CreateEntryParams {
		AccountID: 	entry1.ToAccountID,
		Amount: 	entry1.Amount,
	}

	// create entry 
	entry, err := testQueries.CreateEntry(context.Background(), argTo)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry1.Amount, entry.Amount)
	require.Equal(t, entry1.ToAccountID, entry.AccountID)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)


	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}