package credentials

const (
	// BasicAuthMethod credentials used for basic authentication
	BasicAuthMethod = "basic"
	// KeyFileAuthMethod data used to authenticate through private key file on git
	KeyFileAuthMethod = "keyfile"
	// SSHAgentAuthMethod data used to authenticate through ssh-agent
	SSHAgentAuthMethod = "ssh-agent"
)
