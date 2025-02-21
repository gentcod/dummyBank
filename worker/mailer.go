package worker

import (
	"context"
	"time"

	"encoding/json"
	"fmt"

	db "github.com/gentcod/DummyBank/internal/database"
	"github.com/gentcod/DummyBank/mailer"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
	expiration          = "15m"
	domain              = "http://dummybank.org"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_query", info.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Create cache
	exp, _ := time.ParseDuration(expiration)
	arg := db.RedisData{
		Username: user.Username,
		Email:    user.Email,
	}
	data, err := processor.store.CreateVerifyEmailCache(ctx, arg, exp)
	if err != nil {
		return fmt.Errorf("failed to create verify user email cache: %w", err)
	}

	// SEND EMAIL
	verifyLink := fmt.Sprintf("%s?id=%v&code=%d", domain, data.ID, data.SecretCode)
	recipient := mailer.Recipient{
		Name:             user.Username,
		Email:            user.Email,
		VerificationLink: verifyLink,
	}

	log.Info().Any("%v",data)
	log.Info().Any("%v+",processor.mailer)

	if err := processor.mailer.SendEmail(recipient); err != nil {
		return fmt.Errorf("failed to send user verification email: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg(fmt.Sprintf("processed task: email sent to %s", recipient.Email))

	return nil
}
