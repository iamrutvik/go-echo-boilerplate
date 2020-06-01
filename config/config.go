package config

type (
	// Configurations exported
	Configurations struct {
		Application ApplicationConfigurations
		Server      ServerConfigurations
	}
	ApplicationConfigurations struct {
		Name        string
		JWTExpires  int
		JWTSecret	string
	}
	ServerConfigurations struct {
		Port int
		TLS bool
		CertFile string
		KeyFile string
	}
)

var Settings Configurations
