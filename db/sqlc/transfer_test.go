package db

import (
	"context"
	"database/sql"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createRandomTransfer creates a random amount of transfer by given 2 accounts
func createRandomTransfer(t *testing.T, from Account, to Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        randomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, transfer)

	assert.NotZero(t, transfer.ID)
	assert.NotZero(t, transfer.CreatedAt)

	assert.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	assert.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	assert.Equal(t, arg.Amount, transfer.Amount)

	return transfer
}

// TestCreateTransfer makes sure create a new transfer data in db
func TestCreateTransfer(t *testing.T) {
	transfer := createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))
	assert.NotEmpty(t, transfer)
}

// TestGetTransfer makes sure get transfer
func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, transfer2)

	assert.Equal(t, transfer1.ID, transfer2.ID)
	assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	assert.Equal(t, transfer1.Amount, transfer2.Amount)
	assert.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}

// TestListTransfers makes sure list transfer records with given size and offset
func TestListTransfers(t *testing.T) {
	allTransfers, err := testQueries.ListTransfers(
		context.Background(),
		ListTransfersParams{Limit: math.MaxInt32, Offset: 0},
	)
	assert.NotEmpty(t, allTransfers)
	assert.NoError(t, err)
	total := len(allTransfers)

	count := 10
	for i := 0; i < count; i++ {
		createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))
	}
	transfers, err := testQueries.ListTransfers(
		context.Background(),
		ListTransfersParams{Limit: math.MaxInt32, Offset: 0},
	)
	assert.NotEmpty(t, transfers)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(transfers)-count, total)
}

// TestDeleteTransfer makes sure delete transfer record by given ID
func TestDeleteTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t, createRandomAccount(t), createRandomAccount(t))
	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	assert.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, transfer2)
}
