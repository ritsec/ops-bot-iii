package structs

// MailGunConfig is the config for mailgun
type MailGunConfig struct {
	// Domain is the domain to send emails from
	Domain string

	// APIKey is the API key to use
	APIKey string
}
