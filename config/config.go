package config

type (
	// Configurations exported
	Configurations struct {
		Application ApplicationConfigurations
		Server      ServerConfigurations
	}
	ApplicationConfigurations struct {
		Name        string
		JWTExpireAt int
	}
	ServerConfigurations struct {
		Port int
		TLS bool
		CertFile string
		KeyFile string
	}
)
