package authboss

import (
	"context"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/defaults"
	"gorm.io/gorm"

	_ "github.com/volatiletech/authboss/v3/auth"
	_ "github.com/volatiletech/authboss/v3/recover"
	_ "github.com/volatiletech/authboss/v3/register"
)

// AuthbossService is a service that provides authentication and authorization functionality.
type AuthbossService struct {
	ab     *authboss.Authboss
	db     *gorm.DB
	redis  *redis.Client
	config *Config
}

// NewAuthbossService creates a new AuthbossService.
func NewAuthbossService(db *gorm.DB, redis *redis.Client) (*AuthbossService, error) {
	config := NewConfig()

	// Initialisation of the stores
	serverStorer := NewServerStorer(db)
	sessionStorer := NewSessionStorer(redis)

	// Initialisation of Authboss
	ab := authboss.New()

	// Base config
	ab.Config.Paths.Mount = "/auth"
	ab.Config.Paths.RootURL = config.RootURL

	// Storage
	ab.Config.Storage.Server = serverStorer
	ab.Config.Storage.SessionState = sessionStorer

	// Default config
	defaults.SetCore(&ab.Config, false, false)

	// Initialize
	if err := ab.Init(); err != nil {
		return nil, err
	}

	return &AuthbossService{
		ab:     ab,
		db:     db,
		redis:  redis,
		config: config,
	}, nil
}

// GetAuthboss return the Authboss instance
func (s *AuthbossService) GetAuthboss() *authboss.Authboss {
	return s.ab
}

// GetRouter return the Authboss router
func (s *AuthbossService) GetRouter() http.Handler {
	return s.ab.Config.Core.Router
}

// GetSession returns the session data of a user
func (s *AuthbossService) GetSession(ctx context.Context, key string) (authboss.ClientState, error) {
	return s.ab.Config.Storage.SessionState.Load(ctx, key)
}

// SaveSession saves the session data of a user
func (s *AuthbossService) SaveSession(ctx context.Context, key string, state authboss.ClientState) error {
	return s.ab.Config.Storage.SessionState.Save(ctx, key, state, 24*time.Hour) // Durée configurable
}

// ClearSession clears the session of a user
func (s *AuthbossService) ClearSession(ctx context.Context, key string) error {
	return s.redis.Del(ctx, "authboss_session:"+key).Err()
}

// GetUser returns a user by its ID
func (s *AuthbossService) GetUser(ctx context.Context, id string) (authboss.User, error) {
	return s.ab.Config.Storage.Server.Load(ctx, id)
}
