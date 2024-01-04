package db

import (
	"context"
	"time"
	// "database/sql"
	"testing"

	"github.com/gentcod/DummyBank/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

//createRandomTransfer creates a random Transfer for testing. Returns a Transfer object.
func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTransferParams{
		ID: uuid.New(),
		SenderID: account1.ID,
		RecipientID: account2.ID,
		Amount: util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.SenderID, transfer.SenderID)
	require.Equal(t, arg.RecipientID, transfer.RecipientID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreatetransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGettransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.SenderID, transfer2.SenderID)
	require.Equal(t, transfer1.RecipientID, transfer2.RecipientID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestGetTransfers(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomTransfer(t)
	}
	args := GetTransfersParams{
		Limit: 3,
		Offset: 3,
	}

	transfers, err := testQueries.GetTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 3)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}