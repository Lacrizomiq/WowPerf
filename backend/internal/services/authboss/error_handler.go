package authboss

import (
	"log"
	"net/http"

	"github.com/volatiletech/authboss/v3"
)

// ErrorHandler implements authboss.ErrorHandler for centralized error handling
type ErrorHandler struct {
	renderer *JSONRenderer
}

// NewErrorHandler creates a new ErrorHandler instance
func NewErrorHandler(renderer *JSONRenderer) *ErrorHandler {
	return &ErrorHandler{renderer: renderer}
}

// Wrap implements authboss.ErrorHandler.Wrap
func (h *ErrorHandler) Wrap(handler func(http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		log.Printf("Authboss error: %v", err)

		data := authboss.HTMLData{
			authboss.DataErr: err.Error(),
		}

		// Use the JSONRenderer to output the error
		response, _, renderErr := h.renderer.Render(r.Context(), "error", data)
		if renderErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Determine HTTP status code
		status := determineStatusCode(err)
		w.WriteHeader(status)
		w.Write(response)
	})
}

// determineStatusCode returns the appropriate HTTP status code based on the error
func determineStatusCode(err error) int {
	switch err {
	case authboss.ErrUserNotFound:
		return http.StatusNotFound
	case authboss.ErrTokenNotFound:
		return http.StatusUnauthorized
	case authboss.ErrUserFound:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
