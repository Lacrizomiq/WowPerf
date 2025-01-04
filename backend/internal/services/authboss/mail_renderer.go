package authboss

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/volatiletech/authboss/v3"
)

// MailRenderer implements authboss.Renderer for emails
type MailRenderer struct {
	htmlTemplates map[string]*template.Template
	textTemplates map[string]*template.Template
}

// NewMailRenderer creates a new MailRenderer
func NewMailRenderer() *MailRenderer {
	return &MailRenderer{
		htmlTemplates: make(map[string]*template.Template),
		textTemplates: make(map[string]*template.Template),
	}
}

// Load implements authboss.Renderer.Load
func (m *MailRenderer) Load(names ...string) error {
	for _, name := range names {
		// Load HTML template
		htmlTpl, err := template.New(name + ".html").Parse(defaultHTMLTemplates[name])
		if err != nil {
			return fmt.Errorf("failed to parse HTML template %s: %w", name, err)
		}
		m.htmlTemplates[name] = htmlTpl

		// Load text template
		textTpl, err := template.New(name + ".txt").Parse(defaultTextTemplates[name])
		if err != nil {
			return fmt.Errorf("failed to parse text template %s: %w", name, err)
		}
		m.textTemplates[name] = textTpl
	}
	return nil
}

// Render implements authboss.Renderer.Render
func (m *MailRenderer) Render(ctx context.Context, page string, data authboss.HTMLData) ([]byte, string, error) {
	isHTML := true
	if len(page) > 4 && page[len(page)-4:] == ".txt" {
		isHTML = false
		page = page[:len(page)-4]
	}

	var tpl *template.Template
	var ok bool

	if isHTML {
		tpl, ok = m.htmlTemplates[page]
		if !ok {
			return nil, "", fmt.Errorf("HTML template not found: %s", page)
		}
	} else {
		tpl, ok = m.textTemplates[page]
		if !ok {
			return nil, "", fmt.Errorf("text template not found: %s", page)
		}
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, "", fmt.Errorf("failed to execute template: %w", err)
	}

	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
	}

	return buf.Bytes(), contentType, nil
}

// Default email templates - Modern HTML compatible with email clients
var defaultHTMLTemplates = map[string]string{
	"confirm_email": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Confirm your email</title>
</head>
<body style="margin: 0; padding: 20px; background-color: #f4f4f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;">
    <div style="max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; padding: 20px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);">
        <h1 style="color: #111827; font-size: 24px; margin-bottom: 24px;">Confirm your email address</h1>
        <p style="color: #374151; font-size: 16px; line-height: 24px; margin-bottom: 24px;">
            Click the button below to verify your email address and activate your account.
        </p>
        <div style="text-align: center; margin: 32px 0;">
            <a href="{{.BaseURL}}/auth/confirm?token={{.ConfirmToken}}" 
               style="display: inline-block; background-color: #3b82f6; color: white; text-decoration: none; padding: 12px 24px; border-radius: 6px; font-weight: 500;">
                Confirm Email
            </a>
        </div>
        <p style="color: #6b7280; font-size: 14px; line-height: 20px;">
            If you didn't request this email, you can safely ignore it.
        </p>
        <div style="margin-top: 32px; padding-top: 16px; border-top: 1px solid #e5e7eb;">
            <p style="color: #6b7280; font-size: 12px; line-height: 16px;">
                Button not working? Copy and paste this link into your browser:<br>
                <span style="color: #3b82f6;">{{.BaseURL}}/auth/confirm?token={{.ConfirmToken}}</span>
            </p>
        </div>
    </div>
</body>
</html>`,

	"recover_password": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset your password</title>
</head>
<body style="margin: 0; padding: 20px; background-color: #f4f4f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;">
    <div style="max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; padding: 20px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);">
        <h1 style="color: #111827; font-size: 24px; margin-bottom: 24px;">Reset your password</h1>
        <p style="color: #374151; font-size: 16px; line-height: 24px; margin-bottom: 24px;">
            We received a request to reset your password. Click the button below to choose a new password.
        </p>
        <div style="text-align: center; margin: 32px 0;">
            <a href="{{.BaseURL}}/auth/reset-password?token={{.RecoverToken}}" 
               style="display: inline-block; background-color: #3b82f6; color: white; text-decoration: none; padding: 12px 24px; border-radius: 6px; font-weight: 500;">
                Reset Password
            </a>
        </div>
        <p style="color: #6b7280; font-size: 14px; line-height: 20px;">
            If you didn't request this email, you can safely ignore it. Your password will not be changed.
        </p>
        <div style="margin-top: 32px; padding-top: 16px; border-top: 1px solid #e5e7eb;">
            <p style="color: #6b7280; font-size: 12px; line-height: 16px;">
                Button not working? Copy and paste this link into your browser:<br>
                <span style="color: #3b82f6;">{{.BaseURL}}/auth/reset-password?token={{.RecoverToken}}</span>
            </p>
        </div>
    </div>
</body>
</html>`,
}

// Plain text versions of the templates
var defaultTextTemplates = map[string]string{
	"confirm_email": `
Confirm your email address

Please click the following link to verify your email address:
{{.BaseURL}}/auth/confirm?token={{.ConfirmToken}}

If you didn't request this email, you can safely ignore it.`,

	"recover_password": `
Reset your password

We received a request to reset your password. Click the following link to choose a new password:
{{.BaseURL}}/auth/reset-password?token={{.RecoverToken}}

If you didn't request this email, you can safely ignore it. Your password will not be changed.`,
}
