package configs

// EmailConfig holds the configuration for the email service
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool // TLS bağlantısı kullanımı
	UseSSL   bool // SSL bağlantısı kullanımı
}

// LoadEmailConfig loads the email configuration
func LoadEmailConfig() EmailConfig {
	return EmailConfig{
		Host:     "xxxxxxx",
		Port:     000,
		Username: "xxxxxxxx",
		Password: "xxxxxxxxxx",
		From:     "xxxxxxxxxxx",
		UseTLS:   true,
		UseSSL:   false,
	}
}
