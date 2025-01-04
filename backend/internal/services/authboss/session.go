package authboss

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/volatiletech/authboss/v3"
)

// Default session configuration
const (
	defaultSessionPrefix   = "authboss_session:"
	defaultSessionDuration = 24 * time.Hour
)

// SessionStorer is a storer for the authboss session
type SessionStorer struct {
	redis     *redis.Client
	keyPrefix string
	secure    bool
	domain    string
	sameSite  http.SameSite
}

// SessionConfig holds configuration for the session storer
type SessionConfig struct {
	Redis     *redis.Client
	KeyPrefix string
	Secure    bool
	Domain    string
	SameSite  http.SameSite
}

// NewSessionStorer creates a new SessionStorer
func NewSessionStorer(cfg SessionConfig) *SessionStorer {
	if cfg.KeyPrefix == "" {
		cfg.KeyPrefix = defaultSessionPrefix
	}

	return &SessionStorer{
		redis:     cfg.Redis,
		keyPrefix: cfg.KeyPrefix,
		secure:    cfg.Secure,
		domain:    cfg.Domain,
		sameSite:  cfg.SameSite,
	}
}

// Load implements authboss.ClientStateReadWriter.Load
// It retrieves the session data from Redis using the provided key
func (s *SessionStorer) Load(ctx context.Context, key string) (authboss.ClientState, error) {
	if key == "" {
		return &SessionState{}, nil
	}

	data, err := s.redis.Get(ctx, s.keyPrefix+key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &SessionState{}, nil
		}
		return nil, fmt.Errorf("failed to load session from redis: %w", err)
	}

	var state map[string]string
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	return &SessionState{data: state}, nil
}

// Save implements authboss.ClientStateReadWriter.Save
// It persists the session data to Redis with the specified expiration
func (s *SessionStorer) Save(ctx context.Context, key string, state authboss.ClientState, expire time.Duration) error {
	if state == nil {
		return s.redis.Del(ctx, s.keyPrefix+key).Err()
	}

	// convert to our concrete type
	SessionState, ok := state.(*SessionState)
	if !ok {
		return fmt.Errorf("invalid session state type: %T", state)
	}

	// if there is no data, delete the key
	if len(SessionState.data) == 0 {
		return s.redis.Del(ctx, s.keyPrefix+key).Err()
	}

	// marshal the data to JSON
	data, err := json.Marshal(SessionState.data)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	if expire == 0 {
		expire = defaultSessionDuration
	}

	return s.redis.Set(ctx, s.keyPrefix+key, data, expire).Err()
}

// ReadState implements additional methods required by ClientStateReadWriter
func (s *SessionStorer) ReadState(r *http.Request) (authboss.ClientState, error) {
	// Get sesssion key from cookie
	cookie, err := r.Cookie(authboss.SessionKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return &SessionState{}, nil
		}
		return nil, err
	}

	return s.Load(r.Context(), cookie.Value)
}

// WriteState implements additional methods required by ClientStateReadWriter
func (s *SessionStorer) WriteState(w http.ResponseWriter, r *http.Request, state authboss.ClientState, events []authboss.ClientStateEvent) error {
	for _, ev := range events {
		switch ev.Kind {
		case authboss.ClientStateEventPut:
			// Set new cookie with the session key
			http.SetCookie(w, &http.Cookie{
				Name:     ev.Key,
				Value:    ev.Value,
				Path:     "/",
				Domain:   s.domain,
				Secure:   s.secure,
				HttpOnly: true,
				SameSite: s.sameSite,
				MaxAge:   int(defaultSessionDuration.Seconds()),
			})

		case authboss.ClientStateEventDel:
			// Delete the cookie
			http.SetCookie(w, &http.Cookie{
				Name:     ev.Key,
				Value:    "",
				Path:     "/",
				Domain:   s.domain,
				Secure:   s.secure,
				HttpOnly: true,
				SameSite: s.sameSite,
				MaxAge:   -1,
			})

		case authboss.ClientStateEventDelAll:
			// Delete all cookies except those in whitelist
			whitelist := make(map[string]bool)
			for _, k := range strings.Split(ev.Key, ",") {
				whitelist[strings.TrimSpace(k)] = true
			}

			for _, cookie := range r.Cookies() {
				if !whitelist[cookie.Name] {
					http.SetCookie(w, &http.Cookie{
						Name:     cookie.Name,
						Value:    "",
						Path:     "/",
						Domain:   s.domain,
						Secure:   s.secure,
						HttpOnly: true,
						SameSite: s.sameSite,
						MaxAge:   -1,
					})
				}
			}
		}
	}

	return nil
}

// SessionState implements authboss.ClientState
type SessionState struct {
	data map[string]string
}

// Get implements authboss.ClientState.Get
func (s *SessionState) Get(key string) (string, bool) {
	if s.data == nil {
		return "", false
	}
	val, ok := s.data[key]
	return val, ok
}

// Set implements authboss.ClientState.Set
func (s *SessionState) Set(key, val string) {
	if s.data == nil {
		s.data = make(map[string]string)
	}
	s.data[key] = val
}

// Delete implements authboss.ClientState.Delete
func (s *SessionState) Delete(key string) {
	if s.data != nil {
		delete(s.data, key)
	}
}
