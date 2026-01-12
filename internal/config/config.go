package config

import (
	"mall/internal/rpc"
	"mall/internal/web"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type PGConfig struct {
	Conn string `env:"PG_CONN" env-required:"true"`
}

type AppConfig struct {
	Environment     string `env:"ENVIRONMENT" env-required:"true"`
	LogLevel        string `env:"LOG_LEVEL" env-default:"DEBUG"`
	PG              PGConfig
	Web             web.WebConfig
	Rpc             rpc.RpcConfig
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
}

func InitConfig() (*AppConfig, error) {
	var cfg AppConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
