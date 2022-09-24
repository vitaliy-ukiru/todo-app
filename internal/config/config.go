package config

import (
	"io"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		User     string `yaml:"user" env:"PG_USER"`
		Password string `yaml:"password" env:"PG_PASSWORD"`
		Database string `yaml:"database" env:"PG_DATABASE"`
		Host     string `yaml:"host" env:"PG_HOST"`
		Port     int    `yaml:"port" env:"PG_PORT"`
	} `yaml:"database"`
	Server struct {
		Port     int    `yaml:"port" env-default:"80"`
		Addr     string `yaml:"addr"`
		BasePath string `yaml:"base_path"`
		SSL      *struct {
			CertFile string `yaml:"cert_file"`
			KeyFile  string `yaml:"key_file"`
		} `yaml:"ssl"`
	} `yaml:"server"`
	JWT struct {
		PrivateKeyPath  string        `yaml:"private_key_path" env:"JWT_PRIVATE_KEY"`
		AccessTokenLife time.Duration `yaml:"access_token_life"`
	} `yaml:"jwt"`
}

func New(file io.Reader) (*Config, error) {
	cfg := new(Config)
	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
