package authboss

import (
	"fmt"
	"net/http"

	"github.com/volatiletech/authboss/v3"
)

// Responder implements authboss.HTTPResponder for JSON responses
type Responder struct {
	renderer *JSONRenderer
}

// NewResponder creates a new Responder instance
func NewResponder(renderer *JSONRenderer) *Responder {
	return &Responder{
		renderer: renderer,
	}
}

// Respond implements authboss.HTTPResponder.Respond
func (r *Responder) Respond(w http.ResponseWriter, req *http.Request, code int, page string, data authboss.HTMLData) error {
	// if data is nil, initialize it
	if data == nil {
		data = authboss.HTMLData{}
	}

	// Get any flash messages
	if success := authboss.FlashSuccess(w, req); len(success) > 0 {
		data["flash_success"] = success
	}
	if err := authboss.FlashError(w, req); len(err) > 0 {
		data[authboss.DataErr] = err
	}

	// Add status code to the data
	data["status_code"] = code

	// Render the response using the JSON renderer
	body, contentType, err := r.renderer.Render(req.Context(), page, data)
	if err != nil {
		return fmt.Errorf("error rendering response: %w", err)
	}

	// Set response headers
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)

	// Write the response body
	_, err = w.Write(body)
	return err
}
