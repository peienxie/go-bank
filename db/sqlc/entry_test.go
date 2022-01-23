package db

import (
	"context"
	"database/sql"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createRandomEntry creates a random amount of entry by given account
func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    randomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, entry)

	assert.NotZero(t, entry.ID)
	assert.NotZero(t, entry.CreatedAt)

	assert.Equal(t, arg.AccountID, entry.AccountID)
	assert.Equal(t, arg.Amount, entry.Amount)

	return entry
}

// TestCreateEntry makes sure create a new entry data in db
func TestCreateEntry(t *testing.T) {
	entry := createRandomEntry(t, createRandomAccount(t))
	assert.NotEmpty(t, entry)
}

// TestGetEntry makes sure get entry
func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t, createRandomAccount(t))
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, entry2)

	assert.Equal(t, entry1.ID, entry2.ID)
	assert.Equal(t, entry1.AccountID, entry2.AccountID)
	assert.Equal(t, entry1.Amount, entry2.Amount)
	assert.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
}

// TestListEntries makes sure list entry records with given size and offset
func TestListEntries(t *testing.T) {
	allEntries, err := testQueries.ListEntries(
		context.Background(),
		ListEntriesParams{Limit: math.MaxInt32, Offset: 0},
	)
	assert.NoError(t, err)
	total := len(allEntries)

	var newEntries []Entry
	for i := 0; i < 10; i++ {
		newEntries = append(
			newEntries,
			createRandomEntry(t, createRandomAccount(t)),
		)
	}

	offset := total
	length := 10
	entries, err := testQueries.ListEntries(
		context.Background(),
		ListEntriesParams{Limit: int32(length), Offset: int32(offset)},
	)
	assert.NotEmpty(t, entries)
	assert.NoError(t, err)
	assert.Len(t, entries, length)
	for i := range entries {
		assert.Equal(t, entries[i].ID, newEntries[i].ID)
		assert.Equal(t, entries[i].AccountID, newEntries[i].AccountID)
		assert.Equal(t, entries[i].Amount, newEntries[i].Amount)
		assert.Equal(t, entries[i].CreatedAt, newEntries[i].CreatedAt)
	}
}

// TestDeleteEntry makes sure delete entry record by given ID
func TestDeleteEntry(t *testing.T) {
	entry1 := createRandomEntry(t, createRandomAccount(t))
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	assert.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, entry2)
}
