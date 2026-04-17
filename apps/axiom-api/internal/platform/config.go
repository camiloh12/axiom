package platform

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port        string `envconfig:"PORT" default:"8080"`
	DatabaseURL string `envconfig:"DATABASE_URL" default:"postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable"`
	JWTPrivKey  string `envconfig:"JWT_PRIVATE_KEY"`
	JWTPubKey   string `envconfig:"JWT_PUBLIC_KEY"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
