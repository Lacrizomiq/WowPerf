// services/blizzard/types/protected.go

package types

import (
	"context"
)

// ProtectedClientInterface defines the interface for protected API calls
type ProtectedClientInterface interface {
	MakeProtectedRequest(ctx context.Context, userID uint, endpoint, namespace, locale string) ([]byte, error)
}

// ProfileServiceParams contains the required parameters for the API
type ProfileServiceParams struct {
	Region    string // REQUIRED: us, eu, kr, tw
	Namespace string // REQUIRED: profile-{region}
	Locale    string // The locale to reflect in localized data
}
