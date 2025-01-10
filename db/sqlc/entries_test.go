package db

import (
	"context"
	"github.com/auronvila/simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomEntry(t *testing.T) Entry {
	account1 := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, entry.AccountID, account1.ID)
	return entry
}

func TestCreateEntry(t *testing.T) {
	createdEntry := createRandomEntry(t)
	require.NotEmpty(t, createdEntry)
}

func TestGetEntry(t *testing.T) {
	createdEntry := createRandomEntry(t)
	require.NotEmpty(t, createdEntry)
	getCreatedEntry, err := testQueries.GetEntry(context.Background(), createdEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getCreatedEntry)
	require.Equal(t, createdEntry.ID, getCreatedEntry.ID)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, len(entries))

	for _, entry := range entries {
		_, err := testQueries.GetEntry(context.Background(), entry.ID)
		require.NoError(t, err)
	}
}
