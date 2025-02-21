package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store provides all functions to execute db SQL queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTXResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)

	// CreateVerifyEmailCache stores some data with expiration in redis cache
	CreateVerifyEmailCache(ctx context.Context, arg RedisData, expiration time.Duration) (RedisData, error)

	// GetVerifyEmailCache retrieves stored data in redis cache
	GetVerifyEmailCache(ctx context.Context, key string) (RedisData, error)
}

// SQLStore provides all functions to execute db SQL queries
type SQLStore struct {
	*Queries
	rdb *redis.Client
	db  *sql.DB
}

func NewStore(db *sql.DB, client *redis.Client) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
		rdb:     client,
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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
