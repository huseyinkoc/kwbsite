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
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: "camsihfm@gmail.com",
		Password: "hjaajpuhxswysgzl",
		From:     "camsihfm@gmail.com",
		UseTLS:   true,
		UseSSL:   false,
	}
}
