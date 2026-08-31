package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Manages InterviewState persistence in Redis.

type Repository interface {
	SaveState(
		ctx context.Context,
		state *InterviewState,
	) error

	GetState(
		ctx context.Context,
		interviewID string,
	) (*InterviewState, error)

	DeleteState(
		ctx context.Context,
		interviewID string,
	) error
}

type repository struct {
	redis *redis.Client
}

func NewRepository(redisClient *redis.Client) Repository {
	return &repository{
		redis: redisClient,
	}
}

const interviewStateKeyPrefix = "interview:state:"

// ----------------------------------------------------------------------
// Save Interview State
// ----------------------------------------------------------------------

func (r *repository) SaveState(
	ctx context.Context,
	state *InterviewState,
) error {

	payload, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal interview state: %w", err)
	}

	ttl := time.Until(state.ExpiresAt)
	if ttl <= 0 {
		ttl = time.Hour
	}

	key := interviewStateKeyPrefix + state.InterviewID

	if err := r.redis.Set(ctx, key, payload, ttl).Err(); err != nil {
		return fmt.Errorf("save interview state: %w", err)
	}

	return nil
}

// ----------------------------------------------------------------------
// Get Interview State
// ----------------------------------------------------------------------

func (r *repository) GetState(
	ctx context.Context,
	interviewID string,
) (*InterviewState, error) {

	key := interviewStateKeyPrefix + interviewID

	payload, err := r.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("interview state not found")
		}
		return nil, fmt.Errorf("get interview state: %w", err)
	}

	var state InterviewState

	if err := json.Unmarshal(payload, &state); err != nil {
		return nil, fmt.Errorf("unmarshal interview state: %w", err)
	}

	return &state, nil
}

// ----------------------------------------------------------------------
// Delete Interview State
// ----------------------------------------------------------------------

func (r *repository) DeleteState(
	ctx context.Context,
	interviewID string,
) error {

	key := interviewStateKeyPrefix + interviewID

	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete interview state: %w", err)
	}

	return nil
}