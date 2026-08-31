package interview

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
	"github.com/redis/go-redis/v9"
)

// ----------------------------------------------------------------------
// Session Manager Interface
// ----------------------------------------------------------------------

type SessionManager interface {
	CreateSession(
		ctx context.Context,
		state *engine.InterviewState,
	) error

	GetSession(
		ctx context.Context,
		interviewID string,
	) (*engine.InterviewState, error)

	UpdateSession(
		ctx context.Context,
		state *engine.InterviewState,
	) error

	DeleteSession(
		ctx context.Context,
		interviewID string,
	) error
}

type sessionManager struct {
	redis *redis.Client
	ttl   time.Duration
}

func NewSessionManager(
	redisClient *redis.Client,
	sessionTTL time.Duration,
) SessionManager {

	if sessionTTL == 0 {
		sessionTTL = 2 * time.Hour
	}

	return &sessionManager{
		redis: redisClient,
		ttl:   sessionTTL,
	}
}

// ----------------------------------------------------------------------
// Create Session
// ----------------------------------------------------------------------

func (s *sessionManager) CreateSession(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	key := InterviewSessionKeyPrefix + state.InterviewID

	payload, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal interview state: %w", err)
	}

	if err := s.redis.Set(ctx, key, payload, s.ttl).Err(); err != nil {
		return fmt.Errorf("create interview session: %w", err)
	}

	return nil
}

// ----------------------------------------------------------------------
// Get Session
// ----------------------------------------------------------------------

func (s *sessionManager) GetSession(
	ctx context.Context,
	interviewID string,
) (*engine.InterviewState, error) {

	key := InterviewSessionKeyPrefix + interviewID

	payload, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {

		if err == redis.Nil {
			return nil, ErrInterviewSessionNotFound
		}

		return nil, fmt.Errorf("get interview session: %w", err)
	}

	var state engine.InterviewState

	if err := json.Unmarshal(payload, &state); err != nil {
		return nil, fmt.Errorf("unmarshal interview state: %w", err)
	}

	return &state, nil
}

// ----------------------------------------------------------------------
// Update Session
// ----------------------------------------------------------------------

func (s *sessionManager) UpdateSession(
	ctx context.Context,
	state *engine.InterviewState,
) error {

	key := InterviewSessionKeyPrefix + state.InterviewID

	payload, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal interview state: %w", err)
	}

	if err := s.redis.Set(ctx, key, payload, s.ttl).Err(); err != nil {
		return fmt.Errorf("update interview session: %w", err)
	}

	return nil
}

// ----------------------------------------------------------------------
// Delete Session
// ----------------------------------------------------------------------

func (s *sessionManager) DeleteSession(
	ctx context.Context,
	interviewID string,
) error {

	key := InterviewSessionKeyPrefix + interviewID

	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete interview session: %w", err)
	}

	return nil
}