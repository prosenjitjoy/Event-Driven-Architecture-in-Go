package clickhouse

import "fmt"

type CHConfig struct {
	Host     string `env:"CH_HOST" env-required:"true"`
	Port     string `env:"CH_PORT" env-required:"true"`
	Database string `env:"CH_DB" env-required:"true"`
	Username string `env:"CH_USER" env-required:"true"`
	Password string `env:"CH_PASS" env-required:"true"`
}

func (c CHConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
