package db

import (
	"context"
	"database/sql"
	"github.com/iqbalalfikri/simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := util.RandomMoney()

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}

	result, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	id, err := result.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)

	return Transfer{
		ID:            int32(id),
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	expected := createRandomTransfer(t)

	result, err := testQueries.GetTransfer(context.Background(), expected.ID)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.FromAccountID, result.FromAccountID)
	require.Equal(t, expected.ToAccountID, result.ToAccountID)
	require.Equal(t, expected.Amount, result.Amount)
	require.NotZero(t, result.CreatedAt)
}

func TestUpdateTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)
	newAmount := util.RandomMoney()

	arg := UpdateTransferParams{
		Amount: newAmount,
		ID:     transfer.ID,
	}

	err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)

	result, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, transfer.ID, result.ID)
	require.Equal(t, transfer.FromAccountID, result.FromAccountID)
	require.Equal(t, transfer.ToAccountID, result.ToAccountID)
	require.Equal(t, newAmount, result.Amount)
	require.NotZero(t, result.CreatedAt)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	result, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, result)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
