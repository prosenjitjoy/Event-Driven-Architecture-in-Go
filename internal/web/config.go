package web

import "fmt"

type WebConfig struct {
	Host string `env:"WEB_HOST" env-default:"0.0.0.0"`
	Port string `env:"WEB_PORT" env-default:"8080"`
}

func (c WebConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
