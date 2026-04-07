package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"db_name"`
	Ssl      string `yaml:"ssl"`
}

type CacheConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	DB   int    `yaml:"db"`
}

type Config struct {
	Server   ServerConfig   `yaml:"http_server"`
	Database DatabaseConfig `yaml:"postgres"`
	Cache    CacheConfig    `yaml:"redis"`
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}

	raw, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("error parse config. error %w", err)
	}

	if err := yaml.Unmarshal(raw, cfg); err != nil {
		return nil, fmt.Errorf("error parse config. error %w", err)
	}

	return cfg, nil
}

func (cfg *DatabaseConfig) Dsn() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s name=%s sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Ssl,
	)
}
