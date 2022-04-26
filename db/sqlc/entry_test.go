package db

import (
	"context"
	"database/sql"
	"github.com/iqbalalfikri/simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)

	amount := util.RandomMoney()
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    amount,
	}

	result, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	id, err := result.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)

	return Entry{
		ID:        int32(id),
		AccountID: account.ID,
		Amount:    amount,
	}
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	expected := createRandomEntry(t)

	result, err := testQueries.GetEntry(context.Background(), expected.ID)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.AccountID, result.AccountID)
	require.Equal(t, expected.Amount, result.Amount)
	require.NotZero(t, result.CreatedAt)
}

func TestUpdateEntry(t *testing.T) {
	entry := createRandomEntry(t)
	newAmount := util.RandomMoney()

	arg := UpdateEntryParams{
		Amount: newAmount,
		ID:     entry.ID,
	}

	err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)

	result, err := testQueries.GetEntry(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, entry.ID, result.ID)
	require.Equal(t, entry.AccountID, result.AccountID)
	require.Equal(t, newAmount, result.Amount)
	require.NotZero(t, result.CreatedAt)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	result, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, result)
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
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
