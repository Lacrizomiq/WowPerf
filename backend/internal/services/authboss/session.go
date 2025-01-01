package authboss

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/volatiletech/authboss/v3"
)

// SessionStorer is a storer for the authboss session
type SessionStorer struct {
	redis     *redis.Client
	keyPrefix string
}

// NewSessionStorer creates a new SessionStorer
func NewSessionStorer(redis *redis.Client) *SessionStorer {
	return &SessionStorer{
		redis:     redis,
		keyPrefix: "authboss_session:",
	}
}

// Load implements authboss.ClientStateReadWriter
func (s *SessionStorer) Load(ctx context.Context, key string) (authboss.ClientState, error) {
	data, err := s.redis.Get(ctx, s.keyPrefix+key).Result()
	if err == redis.Nil {
		return nil, authboss.ErrClientStateNotFound
	}
	if err != nil {
		return nil, err
	}

	var state map[string]string
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return nil, err
	}

	return &SessionState{state}, nil
}

// Save implements authboss.ClientStateReadWriter
func (s *SessionStorer) Save(ctx context.Context, key string, state authboss.ClientState, expire time.Duration) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, s.keyPrefix+key, data, expire).Err()
}

// SessionState is a state for the authboss session
type SessionState struct {
	data map[string]string
}

func (s *SessionState) Get(key string) (string, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *SessionState) Set(key, val string) {
	if s.data == nil {
		s.data = make(map[string]string)
	}
	s.data[key] = val
}

func (s *SessionState) Delete(key string) {
	delete(s.data, key)
}
