package config

import (
	"os"
)

type AppEnv struct {
	Host string
	Port string
}

type PostgresEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type EnvConfig struct {
	App    *AppEnv
	GrpcDB *PostgresEnv
}

func NewEnvConfig() *EnvConfig {
	envConfig := &EnvConfig{
		App: &AppEnv{
			Host: os.Getenv("GATEWAY_APP_HOST"),
			Port: os.Getenv("SERVICES_PORT"),
		},
		GrpcDB: &PostgresEnv{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: "grpc_db",
		},
	}
	return envConfig
}
