package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port             int      `envconfig:"PORT" default:"8080"`
	DatabaseURL      string   `envconfig:"DATABASE_URL" required:"true"`
	AllowedOrigins   []string `envconfig:"ALLOWED_ORIGINS" default:"http://localhost:5173"`
	SessionSecret    string   `envconfig:"SESSION_SECRET" required:"true"`
	WebAuthnRPID     string   `envconfig:"WEBAUTHN_RP_ID" default:"localhost"`
	WebAuthnRPName   string   `envconfig:"WEBAUTHN_RP_NAME" default:"StanzaBonanza"`
	WebAuthnOrigins  []string `envconfig:"WEBAUTHN_RP_ORIGINS" default:"http://localhost:5173"`
	ResendAPIKey     string   `envconfig:"RESEND_API_KEY"`
	MagicLinkBaseURL string   `envconfig:"MAGIC_LINK_BASE_URL" default:"http://localhost:5173/auth/verify"`
}

func Load() (*Config, error) {
	// Best-effort load of .env file; real env vars take precedence
	_ = godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
