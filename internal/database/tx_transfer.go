package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TransferTxParams contains the necessary input parameters of the transfer transaction
type TransferTxParams struct {
	SenderID    uuid.UUID `json:"sender_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Amount      int64     `json:"amount"`
}

// TransferTXResult contains the result of the transfer transaction
type TransferTXResult struct {
	Transfer         Transfer `json:"transfer"`
	SenderAccount    Account  `json:"sender_account"`
	RecipientAccount Account  `json:"recipient_account"`
	SenderEntry      Entry    `json:"sender_entry"`
	RecipientEntry   Entry    `json:"recipient_entry"`
}

// var txKey = struct{}{}

// TransferTx performs money transfer from one account to the other.
// It creates a transfer record, adds account entries and update accounts' balance witthin a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTXResult, error) {
	var result TransferTXResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(txKey)
		// fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create sender entry")
		result.SenderEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.SenderID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create recipient entry")
		result.RecipientEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.RecipientID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//To handle deadlock: Sequentially update accounts
		_, senderIntValue := arg.SenderID.Time().UnixTime()
		_, recipientIntValue := arg.RecipientID.Time().UnixTime()

		if senderIntValue < recipientIntValue {
			result.SenderAccount, result.RecipientAccount, err = addMoney(ctx, q, arg.SenderID, -arg.Amount, arg.RecipientID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.RecipientAccount, result.SenderAccount, err = addMoney(ctx, q, arg.RecipientID, arg.Amount, arg.SenderID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

// addMoney performs money transaction from a sender account to a recipients account.
// Returns two account objects and an error object.
func addMoney(ctx context.Context,
	q *Queries,
	account1ID uuid.UUID,
	amount1 int64,
	account2ID uuid.UUID,
	amount2 int64) (account1 Account, account2 Account, err error) {
	senderAccount, err := q.GetAccountForUpdate(ctx, account1ID)
	if err != nil {
		return
	}

	account1, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:        senderAccount.ID,
		Balance:   senderAccount.Balance + amount1,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return
	}

	recipientAccount, err := q.GetAccountForUpdate(ctx, account2ID)
	if err != nil {
		return
	}

	account2, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:        recipientAccount.ID,
		Balance:   recipientAccount.Balance + amount2,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return
	}
	return
}
