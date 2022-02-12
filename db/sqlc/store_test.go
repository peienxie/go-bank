package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTransferTx makes sure money transfer from one account to the other account
func TestTransferTx(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	fmt.Printf("before: from: %d, to: %d\n", fromAccount.Balance, toAccount.Balance)

	n := 10
	amount := int64(10)
	errChan := make(chan error)
	resultChan := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func(i int) {
			result, err := testStore.TransferTx(
				context.Background(),
				TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        amount,
				})
			errChan <- err
			resultChan <- result
		}(i)
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errChan
		assert.NoError(t, err)
		result := <-resultChan
		assert.NotEmpty(t, result)
		fmt.Printf("tx: from: %d, to: %d\n", result.FromAccount.Balance, result.ToAccount.Balance)

		assert.NotEmpty(t, result.Transfer)
		assert.NotZero(t, result.Transfer.ID)
		assert.Equal(t, result.Transfer.FromAccountID, fromAccount.ID)
		assert.Equal(t, result.Transfer.ToAccountID, toAccount.ID)
		assert.Equal(t, result.Transfer.Amount, amount)
		assert.NotZero(t, result.Transfer.CreatedAt)
		_, err = testStore.GetTransfer(context.Background(), result.Transfer.ID)
		assert.NoError(t, err)

		assert.NotEmpty(t, result.FromEntry)
		assert.NotZero(t, result.FromEntry.ID)
		assert.Equal(t, result.FromEntry.AccountID, fromAccount.ID)
		assert.Equal(t, result.FromEntry.Amount, -amount)
		assert.NotZero(t, result.FromEntry.CreatedAt)
		_, err = testStore.GetEntry(context.Background(), result.FromEntry.ID)
		assert.NoError(t, err)

		assert.NotEmpty(t, result.ToEntry)
		assert.NotZero(t, result.ToEntry.ID)
		assert.Equal(t, result.ToEntry.AccountID, toAccount.ID)
		assert.Equal(t, result.ToEntry.Amount, amount)
		assert.NotZero(t, result.ToEntry.CreatedAt)
		_, err = testStore.GetEntry(context.Background(), result.ToEntry.ID)
		assert.NoError(t, err)

		assert.NotEmpty(t, result.FromAccount)
		assert.NotZero(t, result.FromAccount.ID)

		assert.NotEmpty(t, result.ToAccount)
		assert.NotZero(t, result.ToAccount.ID)

		// check balance
		diff1 := fromAccount.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - toAccount.Balance

		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1 > 0)
		assert.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		assert.True(t, k >= 1 && k <= n)
		assert.NotContains(t, existed, k)
	}

	updatedFromAccount, err := testStore.GetAccount(context.Background(), fromAccount.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedFromAccount)
	assert.Equal(t, fromAccount.Balance-amount*int64(n), updatedFromAccount.Balance)

	updatedToAccount, err := testStore.GetAccount(context.Background(), toAccount.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedToAccount)
	assert.Equal(t, toAccount.Balance+amount*int64(n), updatedToAccount.Balance)
	fmt.Printf("After: from: %d, to: %d\n", updatedFromAccount.Balance, updatedToAccount.Balance)
}
