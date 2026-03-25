package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"time"
)

type Config struct {
	HostAddress string        `env:"REST_HOST_ADDRESS"`
	GrpcPort    string        `env:"GRPC_PORT"`
	ServTimeout time.Duration `env:"SERV_TIMEOUT"`
	DBHost      string        `env:"DB_HOST"`
	DBPort      string        `env:"DB_PORT"`
	DBUser      string        `env:"DB_USER"`
	DBPassword  string        `env:"DB_PASSWORD"`
	DBName      string        `env:"DB_NAME"`
	LogLevel    string        `env:"LOG_LEVEL"`
}

func InitConfig() (*Config, error) {
	conf := Config{}

	err := godotenv.Load("internal/config/config.env")
	if err != nil {
		return nil, err
	}

	err = env.Parse(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
