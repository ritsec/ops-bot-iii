package structs

// WebConfig is the config for the web server
type WebConfig struct {
	// Port is the port to listen on
	Port string

	// Hostname is the hostname to listen on
	Hostname string

	// Protocol is the protocol to use
	Protocol string
}
