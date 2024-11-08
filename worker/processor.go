package worker

import (
	"context"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/redis/go-redis/v9"
)

const (
	QueueCritical = "critical"
	QueueDefault = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)
	
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault: 5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().AnErr("failed to process task:", err).Bytes("payload", task.Payload())
			}),
			Logger: logger,
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store: store,
	}
}