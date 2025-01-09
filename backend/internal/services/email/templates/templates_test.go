package templates

import (
	"os"
	"strings"
	"testing"
)

func TestNewTemplateManager(t *testing.T) {
	// Créer un template temporaire pour les tests
	tempDir := t.TempDir()
	templateContent := `<!DOCTYPE html>
<html>
<body>
    <h1>Hello {{.Username}}</h1>
    <p>Click here to reset your password: {{.ResetURL}}</p>
</body>
</html>`

	// Écrire le template dans un fichier temporaire
	err := os.WriteFile(tempDir+"/password_reset.html", []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Modifier temporairement le chemin du template pour les tests
	originalTemplate := emailTemplates[PasswordReset]
	emailTemplates[PasswordReset] = Template{
		Name:    PasswordReset,
		Subject: "Reset Your Password - WowPerf",
		Path:    tempDir + "/password_reset.html",
	}
	defer func() {
		emailTemplates[PasswordReset] = originalTemplate
	}()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Valid template manager initialization",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm, err := NewTemplateManager()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTemplateManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tm == nil {
				t.Error("Expected TemplateManager to not be nil")
			}
		})
	}
}

func TestTemplateManager_RenderTemplate(t *testing.T) {
	// Créer un template temporaire pour les tests
	tempDir := t.TempDir()
	templateContent := `Hello {{.Username}}, reset your password here: {{.ResetURL}}`

	err := os.WriteFile(tempDir+"/password_reset.html", []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Modifier temporairement le chemin du template pour les tests
	originalTemplate := emailTemplates[PasswordReset]
	emailTemplates[PasswordReset] = Template{
		Name:    PasswordReset,
		Subject: "Reset Your Password - WowPerf",
		Path:    tempDir + "/password_reset.html",
	}
	defer func() {
		emailTemplates[PasswordReset] = originalTemplate
	}()

	tests := []struct {
		name       string
		templateID string
		data       interface{}
		want       string
		wantErr    bool
	}{
		{
			name:       "Valid template rendering",
			templateID: PasswordReset,
			data: struct {
				Username string
				ResetURL string
			}{
				Username: "testuser",
				ResetURL: "https://test.com/reset",
			},
			want:    "Hello testuser, reset your password here: https://test.com/reset",
			wantErr: false,
		},
		{
			name:       "Invalid template name",
			templateID: "nonexistent",
			data:       nil,
			want:       "",
			wantErr:    true,
		},
	}

	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create TemplateManager: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tm.RenderTemplate(tt.templateID, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateManager.RenderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.Contains(got, tt.want) {
				t.Errorf("TemplateManager.RenderTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateManager_GetTemplateSubject(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create TemplateManager: %v", err)
	}

	tests := []struct {
		name       string
		templateID string
		want       string
		wantErr    bool
	}{
		{
			name:       "Valid template subject",
			templateID: PasswordReset,
			want:       "Reset Your Password - WowPerf",
			wantErr:    false,
		},
		{
			name:       "Invalid template name",
			templateID: "nonexistent",
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tm.GetTemplateSubject(tt.templateID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateManager.GetTemplateSubject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TemplateManager.GetTemplateSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}
