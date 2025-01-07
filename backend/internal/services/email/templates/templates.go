package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"
)

// Template names constants
const (
	PasswordReset = "password-reset"
)

// Template represents an email template with its metadata
type Template struct {
	Name    string
	Subject string
	Path    string
}

// TemplateData contains the base data available to all templates
type TemplateData struct {
	AppName    string
	AppURL     string
	CustomData interface{}
}

// TemplateManager handles email template loading and rendering
type TemplateManager struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
}

// Available templates
var emailTemplates = map[string]Template{
	PasswordReset: {
		Name:    PasswordReset,
		Subject: "Reset Your Password - WowPerf",
		Path:    "internal/services/email/templates/password_reset.html",
	},
}

// NewTemplateManager creates and initializes a new TemplateManager
func NewTemplateManager() (*TemplateManager, error) {
	tm := &TemplateManager{
		templates: make(map[string]*template.Template),
	}

	// Load all templates at initialization
	if err := tm.loadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return tm, nil
}

// loadTemplates loads all email templates into memory
func (tm *TemplateManager) loadTemplates() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for name, tmpl := range emailTemplates {
		t, err := template.ParseFiles(tmpl.Path)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}
		tm.templates[name] = t
	}
	return nil
}

// RenderTemplate renders an email template with the given data
func (tm *TemplateManager) RenderTemplate(name string, data interface{}) (string, error) {
	tm.mu.RLock()
	tmpl, exists := tm.templates[name]
	tm.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("template %s not found", name)
	}

	// Create a buffer for template execution
	var buf bytes.Buffer

	// Execute template with provided data
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.String(), nil
}

// GetTemplateSubject returns the subject for a given template
func (tm *TemplateManager) GetTemplateSubject(name string) (string, error) {
	tmpl, exists := emailTemplates[name]
	if !exists {
		return "", fmt.Errorf("template %s not found", name)
	}
	return tmpl.Subject, nil
}
