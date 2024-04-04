package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const defaultConfigPath = "./configs/prod.yaml"

type (
	Config struct {
		HTTPServer   HTTPServer `yaml:"http"`
		StatsService GRPCClient `yaml:"stats-service"`
	}

	HTTPServer struct {
		Host    string        `yaml:"host"`
		Port    int           `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"`
	}

	GRPCClient struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
)

func LoadConfig(configPath string) (*Config, error) {

	if configPath == "" {
		if configPathEnv := os.Getenv("CONFIG_PATH"); configPathEnv != "" {
			configPath = configPathEnv
		} else {
			configPath = defaultConfigPath
		}
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
