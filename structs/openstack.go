package structs

// OpenstackConfig is the config for Openstack
type OpenstackConfig struct {
	// Enabled is whether or not openstack self-service is enabled
	Enabled bool

	// MemberID is the ID of the default role Member
	MemberID string

	// CloudsPath is the path to the Clouds.yaml that includes the password
	CloudsPath string
}
