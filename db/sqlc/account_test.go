package db

import (
	"context"
	"database/sql"
	"github.com/iqbalalfikri/simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	owner := util.RandomOwnerName()
	balance := util.RandomMoney()
	currency := util.RandomCurrency()

	arg := CreateAccountParams{
		Owner:    owner,
		Balance:  balance,
		Currency: currency,
	}

	result, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	id, err := result.LastInsertId()
	require.NoError(t, err)
	require.NotZero(t, id)

	return Account{
		ID:       int32(id),
		Owner:    owner,
		Balance:  balance,
		Currency: currency,
	}
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	expected := createRandomAccount(t)

	result, err := testQueries.GetAccount(context.Background(), expected.ID)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.Owner, result.Owner)
	require.Equal(t, expected.Balance, result.Balance)
	require.Equal(t, expected.Currency, result.Currency)
	require.NotZero(t, result.CreatedAt)
}

func TestQueries_UpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	newMoney := util.RandomMoney()
	arg := UpdateAccountParams{
		Balance: newMoney,
		ID:      account.ID,
	}

	err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)

	result, err := testQueries.GetAccount(context.Background(), arg.ID)
	require.NoError(t, err)

	require.NotEmpty(t, result)

	require.Equal(t, account.ID, result.ID)
	require.Equal(t, account.Owner, result.Owner)
	require.Equal(t, newMoney, result.Balance)
	require.Equal(t, account.Currency, result.Currency)

}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	result, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, result)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
