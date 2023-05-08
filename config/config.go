package config

import (
	"os"
	"path/filepath"

	"flag"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		HTTP *HttpServerConfig `yaml:"http"`
		GRPC *GrpcServerConfig `yaml:"grpc"`
	} `yaml:"server"`
	Storages *struct {
		Postgres *DBConfig    `yaml:"postgres"`
		Minio    *MinioConfig `yaml:"minio"`
	}
}

type HttpServerConfig struct {
	BaseURL string `yaml:"base_url"`
	Listen  string `yaml:"listen"`
}

type GrpcServerConfig struct {
	Listen string `yaml:"listen"`
}

type DBConfig struct {
	UseEnv bool   `yaml:"use_env"`
	DSN    string `yaml:"dsn"`
}

type MinioConfig struct {
	URL string `yaml:"url"`
}

func GetConfig() (*Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "Used for set path to config file.")
	flag.Parse()

	var cfg Config
	data, err := os.ReadFile(filepath.Clean("config/config.yaml"))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Storages.Postgres.UseEnv {
		cfg.Storages.Postgres.DSN = os.Getenv("DB_DSN")
	}

	return &cfg, err
}
