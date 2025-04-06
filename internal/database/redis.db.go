package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gentcod/DummyBank/util"
	"github.com/google/uuid"
)

const (
	testExpiration = "5m"
	scMin          = 100000
	scMax          = 999999
	emailKey       = "verify-email:"
)

type RedisData struct {
	ID              uuid.UUID
	Username, Email string
	SecretCode      int64
}

func (store *SQLStore) CreateVerifyEmailCache(ctx context.Context, arg RedisData, expiration time.Duration) (RedisData, error) {
	arg.ID = uuid.New()
	arg.SecretCode = util.RandomInt(scMin, scMax)

	data, err := json.Marshal(arg)
	if err != nil {
		return arg, err
	}

	err = store.rdb.Set(ctx, emailKey+arg.Username, data, expiration).Err()
	if err != nil {
		return arg, err
	}

	return arg, nil
}

func (store *SQLStore) GetVerifyEmailCache(ctx context.Context, key string) (RedisData, error) {
	var result RedisData
	datastr, err := store.rdb.Get(ctx, emailKey+key).Result()
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(datastr), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (store *SQLStore) DeleteVerifyEmailCache(ctx context.Context, key string) error {
	_, err := store.rdb.Del(ctx, emailKey+key).Result()
	if err != nil {
		return err
	}

	return nil
}
