package rpc

import "fmt"

type RpcConfig struct {
	Host string `env:"RPC_HOST" env-default:"0.0.0.0"`
	Port string `env:"RPC_PORT" env-default:"8085"`
}

func (c RpcConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
