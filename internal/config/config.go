package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"role-leader/internal/logger"
	"role-leader/internal/postgres"
)

type Config struct {
	Postgres postgres.Config `yaml:"POSTGRES" env:"POSTGRES"`
	GRPC     GrpcConfig      `yaml:"GRPC" env:"GRPC"`
	REST     RestConfig      `yaml:"REST" env:"REST"`
	Logger   logger.Config   `yaml:"LOGGER" env:"LOGGER"`
}

type GrpcConfig struct {
	Port int `yaml:"ROLE_LEADER_GRPC_PORT" env:"ROLE_LEADER_GRPC_PORT"`
}

type RestConfig struct {
	Port int `yaml:"ROLE_LEADER_REST_PORT" env:"ROLE_LEADER_REST_PORT"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("unable to read config.yaml: %w", err)
	}

	return &cfg, nil
}
