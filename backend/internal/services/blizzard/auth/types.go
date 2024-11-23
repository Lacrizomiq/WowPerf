package auth

import "errors"

var (
	// Auth errors
	ErrTokenExpired       = errors.New("battle.net token expired")
	ErrTokenInvalid       = errors.New("battle.net token invalid")
	ErrNoToken            = errors.New("no battle.net token found")
	ErrAccountNotLinked   = errors.New("battle.net account not linked")
	ErrTokenRefreshFailed = errors.New("failed to refresh battle.net token")
	ErrInvalidScope       = errors.New("insufficient battle.net permissions")
)

// RequiredScopes defines the scopes required for different operations
var (
	ProfileScopes    = []string{"wow.profile"}
	CharacterScopes  = []string{"wow.profile", "profile"}
	CollectionScopes = []string{"wow.profile", "collections"}
)

// BattleNetStatus represents the current status of a Battle.net account link
type BattleNetStatus struct {
	Linked    bool   `json:"linked"`
	BattleTag string `json:"battle_tag,omitempty"` // Changed from *string to string
}

// TokenInfo contains the information of the Battle.net token
type TokenInfo struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	ExpiresIn    int64
	Scope        string
}

// BattleNetProfile contains the information of the Battle.net profile
type BattleNetProfile struct {
	ID        string `json:"id"`
	BattleTag string `json:"battletag"`
	Sub       string `json:"sub"`
}
