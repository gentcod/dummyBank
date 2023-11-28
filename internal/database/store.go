package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

//Provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

//Executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

//Contains the necessary input parameters of the transfer transaction
type TransferTxParams struct {
	SenderID uuid.UUID `json:"sender_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Amount int64 `json:"amount"`
}

//Contains the result of the transfer transaction
type TransferTXResult struct {
	Transfer	Transfer	`json:"transfer"`
	SenderAccount	Account	`json:"sender_account"`
	RecipientAccount	Account	`json:"recipient_account"`
	SenderEntry	Entry	`json:"sender_entry"`
	RecipientEntry	Entry	`json:"recipient_entry"`
}

// var txKey = struct{}{}

//Performs money transfer from one account to the other.
//It creates a transfer record, adds account entries and update accounts' balance witthin a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTXResult, error){
	var result TransferTXResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(txKey)
		// fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			ID: uuid.New(),
			SenderID: arg.SenderID,
			RecipientID: arg.RecipientID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create sender entry")
		result.SenderEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			ID: uuid.New(),
			AccountID: arg.SenderID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create recipient entry")
		result.RecipientEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			ID: uuid.New(),
			AccountID: arg.RecipientID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "get sender account for update")
		senderAccount, err := q.GetAccountForUpdate(ctx, arg.SenderID)
		if err != nil {
			return err
		}

		// fmt.Println(txName, "update sender account")
		result.SenderAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.SenderID,
			Balance: senderAccount.Balance - arg.Amount,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "get sender account for update")
		recipientAccount, err := q.GetAccountForUpdate(ctx, arg.RecipientID)
		if err != nil {
			return err
		}
		// fmt.Println(txName, "update sender account")
		result.RecipientAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.RecipientID,
			Balance: recipientAccount.Balance + arg.Amount,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
