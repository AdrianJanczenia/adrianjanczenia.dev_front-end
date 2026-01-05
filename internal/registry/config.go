package registry

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
	Api struct {
		GatewayURL string
	}
	Templates struct {
		Path string
	}
}

func LoadConfig() (*Config, error) {
	type yamlConfig struct {
		Server struct {
			Port         string `yaml:"port"`
			ReadTimeout  int    `yaml:"readTimeoutSeconds"`
			WriteTimeout int    `yaml:"writeTimeoutSeconds"`
		} `yaml:"server"`
		Api struct {
			GatewayURL string `yaml:"gatewayUrl"`
		} `yaml:"api"`
		Templates struct {
			Path string `yaml:"path"`
		} `yaml:"templates"`
	}

	env := os.Getenv("APP_ENV")
	if env != "production" {
		env = "local"
	}
	configPath := filepath.Join("config", env, "config.yml")
	log.Printf("INFO: loading configuration from %s", configPath)

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var yc yamlConfig
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&yc); err != nil {
		return nil, err
	}

	cfg := &Config{}
	cfg.Server.Port = yc.Server.Port
	cfg.Server.ReadTimeout = time.Duration(yc.Server.ReadTimeout) * time.Second
	cfg.Server.WriteTimeout = time.Duration(yc.Server.WriteTimeout) * time.Second
	cfg.Api.GatewayURL = yc.Api.GatewayURL
	cfg.Templates.Path = yc.Templates.Path

	overrideFromEnv("GATEWAY_URL", &cfg.Api.GatewayURL)

	return cfg, nil
}

func overrideFromEnv(envKey string, configValue *string) {
	if value, exists := os.LookupEnv(envKey); exists && value != "" {
		*configValue = value
	}
}
