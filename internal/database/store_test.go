package db

import (
	"context"
	"testing"

	"github.com/gentcod/DummyBank/util"
	"github.com/stretchr/testify/require"
)

// "context"
// "database/sql"
// "fmt"

// "github.com/google/uuid"

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := util.RandomMoney()

	arg := TransferTxParams{
		SenderID: account1.ID,
		RecipientID: account2.ID,
		Amount: amount,
	}
	
	// run n concurrent transfer transactions
	n := 5

	errs := make(chan error)
	results := make(chan TransferTXResult)

	for i := 0; i < n; i++ {
		go func ()  {
			result, err := store.TransferTx(context.Background(), arg)

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
	
		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)

		require.Equal(t, arg.SenderID, transfer.SenderID)
		require.Equal(t, arg.RecipientID, transfer.RecipientID)
		require.Equal(t, arg.Amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		//check if transfer is added to db
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		senderEntry := result.SenderEntry
		require.NotEmpty(t, senderEntry)
		require.Equal(t, arg.SenderID, senderEntry.AccountID)
		require.Equal(t, -amount, senderEntry.Amount)
		require.NotZero(t, senderEntry.ID)
		require.NotZero(t, senderEntry.CreatedAt)

		//check if sender entry is added to db
		_, err = store.GetEntry(context.Background(), senderEntry.ID)
		require.NoError(t, err)

		recipientEntry := result.RecipientEntry
		require.NotEmpty(t, recipientEntry)
		require.Equal(t, arg.RecipientID, recipientEntry.AccountID)
		require.Equal(t, amount, recipientEntry.Amount)
		require.NotZero(t, recipientEntry.ID)
		require.NotZero(t, recipientEntry.CreatedAt)

		//check if recipient entry is added to db
		_, err = store.GetEntry(context.Background(), recipientEntry.ID)
		require.NoError(t, err)

		//TODO: check accounts balance
	}
}