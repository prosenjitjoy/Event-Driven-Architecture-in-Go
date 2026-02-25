package config

import (
	"mall/internal/clickhouse"
	"mall/internal/rpc"
	"mall/internal/web"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type PGConfig struct {
	Conn string `env:"PG_CONN" env-required:"true"`
}

type NatsConfig struct {
	URL    string `env:"NATS_URL" env-required:"true"`
	Stream string `env:"NATS_STREAM_NAME" env-default:"mall"`
}

type AppConfig struct {
	Environment     string        `env:"ENVIRONMENT" env-required:"true"`
	LogLevel        string        `env:"LOG_LEVEL" env-default:"DEBUG"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
	PG              PGConfig
	CH              clickhouse.CHConfig
	Nats            NatsConfig
	Web             web.WebConfig
	Rpc             rpc.RpcConfig
	ServiceName     string `env:"OTEL_SERVICE_NAME" env-required:"true"`
	rpcServices     string `env:"RPC_SERVICES" env-default:""`
}

func InitConfig() (*AppConfig, error) {
	var cfg AppConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	cfg.Rpc.Services.Decode(cfg.rpcServices)

	return &cfg, nil
}
