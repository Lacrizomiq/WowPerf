// internal/services/email/templates.go
package email

import (
	"bytes"
	"fmt"
	"html/template"
)

// Template names
const (
	TemplateResetPassword = "reset-password"
)

// EmailTemplate represents an email template
type EmailTemplate struct {
	Name    string
	Subject string
	Body    string
}

var templates = map[string]EmailTemplate{
	TemplateResetPassword: {
		Name:    TemplateResetPassword,
		Subject: "Reset Your Password",
		Body: `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your Password</title>
</head>
<body>
    <h2>Password Reset Request</h2>
    <p>Hello,</p>
    <p>You have requested to reset your password. Click the link below to set a new password:</p>
    <p><a href="{{.ResetURL}}">Reset Password</a></p>
    <p>This link will expire in 1 hour.</p>
    <p>If you did not request this password reset, please ignore this email.</p>
    <br>
    <p>Best regards,</p>
    <p>The WowPerf Team</p>
</body>
</html>`,
	},
}

func RenderTemplate(templateName string, data interface{}) (EmailData, error) {
	tmpl, ok := templates[templateName]
	if !ok {
		return EmailData{}, fmt.Errorf("template %s not found", templateName)
	}

	t, err := template.New(tmpl.Name).Parse(tmpl.Body)
	if err != nil {
		return EmailData{}, err
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return EmailData{}, err
	}

	return EmailData{
		Subject: tmpl.Subject,
		HTML:    body.String(),
	}, nil
}
