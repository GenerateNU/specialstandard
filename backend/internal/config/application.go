package config

type Application struct {
	Port           string `env:"PORT, default=8080"`
	Environment    string `env:"ENVIRONMENT, default=development"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS, default=http://localhost:3000"`
}
