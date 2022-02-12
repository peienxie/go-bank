package db

import (
	"context"
	"database/sql"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createRandomAccount generates a random account data for testing
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Username: randomUsername(),
		Balance:  randomMoney(),
		Currency: randomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account)

	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	assert.Equal(t, account.Username, arg.Username)
	assert.Equal(t, account.Balance, arg.Balance)
	assert.Equal(t, account.Currency, arg.Currency)

	return account
}

// TestCreateAccount makes sure create a new account in db
func TestCreateAccount(t *testing.T) {
	account := createRandomAccount(t)
	assert.NotEmpty(t, account)
}

// TestGetAccount makes sure get account by given ID
func TestGetAaccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Username, account2.Username)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

// TestListAccount makes sure list accounts with given size and offset
func TestListAaccount(t *testing.T) {
	allAccounts, err := testQueries.ListAccounts(
		context.Background(),
		ListAccountsParams{Limit: math.MaxInt32, Offset: 0},
	)
	assert.NoError(t, err)
	total := len(allAccounts)

	var newAccounts []Account
	for i := 0; i < 10; i++ {
		newAccounts = append(newAccounts, createRandomAccount(t))
	}

	offset := total
	accounts, err := testQueries.ListAccounts(
		context.Background(),
		ListAccountsParams{Limit: 10, Offset: int32(offset)},
	)
	assert.NotEmpty(t, accounts)
	assert.NoError(t, err)
	assert.Len(t, accounts, 10)
	for i := range accounts {
		assert.Equal(t, accounts[i].ID, newAccounts[i].ID)
		assert.Equal(t, accounts[i].Username, newAccounts[i].Username)
		assert.Equal(t, accounts[i].Balance, newAccounts[i].Balance)
		assert.Equal(t, accounts[i].Currency, newAccounts[i].Currency)
		assert.Equal(t, accounts[i].CreatedAt, newAccounts[i].CreatedAt)
	}
}

// TestUpdateAccountBalance makes sure update account amount of balance by given ID
func TestUpdateAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountBalanceParams{
		ID:      account1.ID,
		Balance: randomMoney(),
	}

	account2, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Username, account2.Username)
	assert.Equal(t, arg.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

// TestAddAccountBalance makes sure add account balance by provided amount
func TestAddAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := AddAccountBalanceParams{
		ID:     account1.ID,
		Amount: randomMoney(),
	}

	account2, err := testQueries.AddAccountBalance(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Username, account2.Username)
	assert.Equal(t, account1.Balance+arg.Amount, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.Equal(t, account1.CreatedAt, account2.CreatedAt)
}

// TestDeleteAccount makes sure delete account record by given ID
func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, account2)
}
