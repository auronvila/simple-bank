package db

import (
	"context"
	"database/sql"
	"github.com/auronvila/simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
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
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
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
}

func TestGetAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), createdAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, createdAccount.ID, account2.ID)
	require.Equal(t, createdAccount.Owner, account2.Owner)
	require.Equal(t, createdAccount.Balance, account2.Balance)
	require.Equal(t, createdAccount.Currency, account2.Currency)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListUserAccount(t *testing.T) {
	owners := []Account{}
	for i := 0; i < 10; i++ {
		owner := createRandomAccount(t)
		owners = append(owners, owner)
	}

	for _, owner := range owners {
		arg := ListAccountsParams{
			Limit:  5,
			Offset: 0,
			Owner:  owner.Owner,
		}
		accounts, err := testQueries.ListAccounts(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, accounts, 1)
	}
}
