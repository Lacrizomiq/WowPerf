package providers

// EmailData is the data structure for sending an email
type EmailData struct {
	To      string
	Subject string
	HTML    string
}

// Provider is an interface that all email providers must implement
type Provider interface {
	// SendEmail sends an email to the given recipient with the provided data
	SendEmail(data EmailData) error

	// Initialize sets up any necessary resources for the provider
	Initialize() error

	// Close cleans up any resources used by the provider
	Close() error
}
