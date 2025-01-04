package authboss

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/volatiletech/authboss/v3"
)

// JSONRenderer implements the authboss.Renderer interface for JSON responses
type JSONRenderer struct{}

// NewJSONRenderer creates a new JSONRenderer
func NewJSONRenderer() *JSONRenderer {
	return &JSONRenderer{}
}

// Load implements the authboss.Renderer
// Since i render json, i don't need to load anything
func (r *JSONRenderer) Load(names ...string) error {
	return nil
}

// Render implements the authboss.Renderer
// Converts the data into a standardized JSON response format
func (r *JSONRenderer) Render(ctx context.Context, page string, data authboss.HTMLData) ([]byte, string, error) {
	response := make(map[string]interface{})

	// Handle validation errors
	if errList, ok := data[authboss.DataValidation]; ok {
		response["errors"] = errList
	}

	// Handle single error message
	if err, ok := data[authboss.DataErr]; ok {
		response["error"] = err
	}

	// Handle preserved data (used for form resubmission)
	if preserve, ok := data[authboss.DataPreserve]; ok {
		response["preserved"] = preserve
	}

	// Handle sucess message (if any)
	if flash, ok := data["flash_success"]; ok {
		response["message"] = flash
	}

	// Add any remaining data
	for k, v := range data {
		if k != authboss.DataValidation &&
			k != authboss.DataErr &&
			k != authboss.DataPreserve &&
			k != "flash_success" {
			response[k] = v
		}
	}

	// add the current page/action
	response["action"] = page

	output, err := json.Marshal(response)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return output, "application/json", nil
}
