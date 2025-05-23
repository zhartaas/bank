package db

import (
	"context"
	"exampleproject/util"
	"github.com/jackc/pgx/v5"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccounts(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccounts(t)
}

func TestGetAccount(t *testing.T) {
	expected := createRandomAccounts(t)
	output, err := testQueries.GetAccount(context.Background(), expected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, output)

	require.Equal(t, expected.ID, output.ID)
	require.Equal(t, expected.Owner, output.Owner)
	require.Equal(t, expected.Balance, output.Balance)
	require.Equal(t, expected.Currency, output.Currency)
	require.WithinDuration(t, expected.CreatedAt.Time, output.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	expected := createRandomAccounts(t)

	expectedBalance := util.RandomMoney()
	expected.Balance = expectedBalance

	arg := UpdateAccountParams{ID: expected.ID, Balance: expectedBalance}
	err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)

	output, err := testQueries.GetAccount(context.Background(), expected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, output)

	require.Equal(t, expected.ID, output.ID)
	require.Equal(t, expected.Owner, output.Owner)
	require.Equal(t, expected.Balance, output.Balance)
	require.Equal(t, expected.Currency, output.Currency)
	require.WithinDuration(t, expected.CreatedAt.Time, output.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	randomAccount := createRandomAccounts(t)

	err := testQueries.DeleteAccount(context.Background(), randomAccount.ID)
	require.NoError(t, err)

	output, err := testQueries.GetAccount(context.Background(), randomAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, output)
}
