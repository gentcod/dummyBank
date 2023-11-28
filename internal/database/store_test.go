package db

import (
	"context"
	"testing"

	"github.com/gentcod/DummyBank/util"
	"github.com/stretchr/testify/require"
)

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
	existed := make(map[int]bool)

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

		//check accounts
		senderAccount := result.SenderAccount
		require.NotEmpty(t, senderAccount)
		require.Equal(t, account1.ID, senderAccount.ID)

		recipientAccount := result.RecipientAccount
		require.NotEmpty(t, recipientAccount)
		require.Equal(t, account2.ID, recipientAccount.ID)

		//check accounts balance
		diff1 := account1.Balance - senderAccount.Balance
		diff2 := recipientAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check the final balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}