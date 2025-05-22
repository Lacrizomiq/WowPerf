package googleauth

import (
	"fmt"
	"wowperf/internal/models"
)

// === TYPES ET STRUCTURES ===

// GoogleUserInfo structures des infos utilisateur récupérées de Google
type GoogleUserInfo struct {
	ID            string `json:"sub"`            // Google ID unique
	Email         string `json:"email"`          // Email utilisateur
	VerifiedEmail bool   `json:"verified_email"` // CRITIQUE : email vérifié par Google
	Name          string `json:"name"`           // Nom complet
	GivenName     string `json:"given_name"`     // Prénom
	FamilyName    string `json:"family_name"`    // Nom de famille
	Picture       string `json:"picture"`        // URL avatar
	Locale        string `json:"locale"`         // Locale (fr, en, etc.)
}

// AuthResult résultat de l'authentification Google
type AuthResult struct {
	User      *models.User `json:"user"`        // Utilisateur authentifié
	IsNewUser bool         `json:"is_new_user"` // true si nouvel utilisateur créé
	Method    string       `json:"method"`      // "login", "signup", ou "link"
}

// UserLookupResult résultat de la recherche utilisateur (logique complexe Google)
type UserLookupResult struct {
	ExistingUser    *models.User `json:"existing_user,omitempty"`
	FoundByGoogleID bool         `json:"found_by_google_id"` // Trouvé par Google ID
	FoundByEmail    bool         `json:"found_by_email"`     // Trouvé par email
	CanAutoLink     bool         `json:"can_auto_link"`      // Peut être lié automatiquement
}

// OAuthError erreur spécifique OAuth
type OAuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *OAuthError) Error() string {
	return fmt.Sprintf("OAuth error [%s]: %s", e.Code, e.Message)
}
