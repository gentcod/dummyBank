package db

import (
	"context"
	"fmt"
	"strconv"
)

type VerfiyEmailParams struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (store *SQLStore) VerfiyEmailTx(ctx context.Context, arg VerfiyEmailParams) (VerifyUserEmailRow, error) {
	var user VerifyUserEmailRow

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		cache, err := store.GetVerifyEmailCache(ctx, arg.ID)
		if err != nil {
			return fmt.Errorf("failed to get cached user data: %v", err)
		}

		secretCode, _ := strconv.Atoi(arg.Token)
		if cache.SecretCode != int64(secretCode) {
			return fmt.Errorf("invalid token")
		}

		user, err = store.VerifyUserEmail(ctx, cache.Email)
		if err != nil {
			return err
		}

		err = store.DeleteVerifyEmailCache(ctx, cache.Username)
		if err != nil {
			return fmt.Errorf("failed to delete cached user data: %v", err)
		}

		return nil
	})

	return user, err
}
