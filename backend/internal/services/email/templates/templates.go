package templates

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"sync"
)

//go:embed password_reset.html
var templateFS embed.FS

// Template names constants
const (
	PasswordReset = "password-reset"
)

// Template represents an email template with its metadata
type Template struct {
	Name    string
	Subject string
	Path    string
	DevPath string // Chemin pour le développement local
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
		Path:    "password_reset.html",                                   // Pour l'embedded FS
		DevPath: "internal/services/email/templates/password_reset.html", // Pour le développement local
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
		// Essayer d'abord l'embedded FS
		content, err := templateFS.ReadFile(tmpl.Path)
		if err == nil {
			t, err := template.New(name).Parse(string(content))
			if err != nil {
				return fmt.Errorf("failed to parse embedded template %s: %w", name, err)
			}
			tm.templates[name] = t
			continue
		}

		// Si l'embedded FS échoue, essayer le chemin de développement
		t, err := template.ParseFiles(tmpl.DevPath)
		if err != nil {
			return fmt.Errorf("failed to parse template %s from either embedded or development path: %w", name, err)
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
